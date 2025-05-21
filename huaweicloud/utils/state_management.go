package utils

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// func RefreshObjectParamOriginValues(objectParamKeys []string) schema.CustomizeDiffFunc {
// 	return func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
// 		var (
// 			err  error
// 			mErr *multierror.Error
// 		)

// 		for _, paramKey := range objectParamKeys {
// 			originParamKey := fmt.Sprintf("%s_origin", paramKey)

// 			// Obtain the origin value.
// 			_, newVal := d.GetChange(originParamKey)

// 			err = d.SetNew(originParamKey, newVal)
// 			if err != nil {
// 				log.Printf("[DEBUG] Unable to set the origin value (corresponding key: %s) because %v, the value is %v",
// 					paramKey, err, newVal)
// 			}

// 			err = d.Clear(originParamKey)
// 			if err != nil {
// 				log.Printf("[DEBUG] Unable to clear the origin value (corresponding key: %s) because %v, the value is %v",
// 					paramKey, err, newVal)
// 			}
// 		}
// 		return mErr.ErrorOrNil()
// 	}
// }

// If the request is successful, obtain the values of all JSON|object parameters first and save them to the
// corresponding '_origin' attributes for subsequent determination and construction of the request body during
// next updates.
// And whether corresponding parameters are changed, the origin values must be refreshed.
func RefreshObjectParamOriginValues(objectParamKeys []string) schema.CustomizeDiffFunc {
	return func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
		var mErr *multierror.Error

		for _, paramKey := range objectParamKeys {
			parts := strings.Split(paramKey, ".")
			log.Printf("[Lance][%s] The parts list is %v", paramKey, parts)
			// Construct the corresponding _origin path.
			originParts := make([]string, len(parts))
			copy(originParts, parts)
			lastIdx := len(originParts) - 1
			originParts[lastIdx] += "_origin"

			// Obtain the origin value.
			rawVal, err := getNestedValue(d, parts)
			if err != nil {
				log.Printf("[DEBUG] failed to get origin value for the parameter '%s': %v", paramKey, err)
				// If the acquisition fails, the subsequent operation of the current parameter is skipped because this
				// parameter may not be configured.
				continue
			}
			log.Printf("[Lance][%s] The rawVal that returned from the function getNestedValue is: %v", paramKey, rawVal)

			// Setting the origin value to the origin attribute.
			if err := setNestedValue(d, originParts, rawVal, true); err != nil {
				mErr = multierror.Append(mErr, fmt.Errorf("failed to set origin value for '%s': %v", paramKey, err))
			}
		}
		return mErr.ErrorOrNil()
	}
}

// func RefreshObjectParamOriginValues(d *schema.ResourceDiff, objectParamKeys []string) error {
// 	var mErr *multierror.Error

// 	for _, key := range objectParamKeys {
// 		parts := strings.Split(key, ".")
// 		// Construct the corresponding _origin path.
// 		originParts := make([]string, len(parts))
// 		copy(originParts, parts)
// 		lastIdx := len(originParts) - 1
// 		originParts[lastIdx] += "_origin"

// 		// Obtain the origin value
// 		rawVal, err := getNestedValue(d, parts)
// 		if err != nil {
// 			log.Printf("[DEBUG] failed to get origin value for the parameter '%s': %v", key, err)
// 			// If the acquisition fails, the subsequent operation of the current parameter is skipped because this
// 			// parameter may not be configured.
// 			continue
// 		}

// 		// Setting the origin value
// 		if err := setNestedValue(d, originParts, rawVal); err != nil {
// 			mErr = multierror.Append(mErr, fmt.Errorf("failed to set origin value for '%s': %v", key, err))
// 		}
// 	}

// 	return mErr.ErrorOrNil()
// }

// getNestedValue method that used to obtain nested values ​​based on the path recursively, because the nested parameter
// must ensure that the complete structure nesting of its corresponding subscript is obtained (only the corresponding
// index is covered)
func getNestedValue(d *schema.ResourceDiff, parts []string) (interface{}, error) {
	var current interface{}
	_, current = d.GetChange(parts[0])
	log.Printf("[Lance] The parts list is: %v", parts)

	for i := 1; i < len(parts); i++ {
		part := parts[i]
		switch cv := current.(type) {
		case []interface{}:
			if len(cv) == 0 {
				return nil, fmt.Errorf("empty list at '%s'", strings.Join(parts[:i+1], "."))
			}
			// Processing lists/collections (automatically taking the first element if the index number is missing).
			current = cv[0]
			if index, err := strconv.Atoi(part); err == nil {
				if index >= len(cv) {
					return nil, fmt.Errorf("index %d out of range", index)
				}
				current = cv[index]
			} else {
				elem, ok := current.(map[string]interface{})
				if !ok {
					return nil, fmt.Errorf("invalid nested path at '%s'", strings.Join(parts[:i+1], "."))
				}
				current = elem[part]
			}
		case map[string]interface{}:
			var ok bool
			current, ok = cv[part]
			if !ok {
				return nil, fmt.Errorf("missing key '%s'", part)
			}
		default:
			return nil, fmt.Errorf("unsupported type at '%s'", strings.Join(parts[:i+1], "."))
		}
	}
	log.Printf("[Lance] The current value is: %v", current)
	return current, nil
}

// setNestedValue method that used to set nested value recursively, because nested parameters must set their full
// structure nesting according to their index (only overwrite the corresponding index).
func setNestedValue(d *schema.ResourceDiff, parts []string, value interface{}, clearAllElems bool) error {
	rootKey := parts[0]
	current := d.Get(rootKey)

	updated, err := updateNestedStructure(current, parts[1:], value)
	if err != nil {
		return err
	}

	log.Printf("[Lance] The updated value is: %v", updated)
	err = d.SetNew(rootKey, updated)
	if err != nil {
		return err
	}
	return nil
}

func updateNestedStructure(current interface{}, parts []string, value interface{}) (interface{}, error) {
	if len(parts) == 0 {
		return value, nil
	}

	part := parts[0]
	switch cv := current.(type) {
	case []interface{}:
		if len(cv) == 0 {
			return nil, errors.New("cannot update empty list")
		}
		// Considering that the index of the Set type is inconsistent during the change before and after, currently only
		// the first element of the List type is automatically processed (applicable to the MaxItems=1 scenario).
		updatedElem, err := updateNestedStructure(cv[0], parts[1:], value)
		if err != nil {
			return nil, err
		}
		cv[0] = updatedElem
		return cv, nil
	case map[string]interface{}:
		subVal, ok := cv[part]
		if !ok {
			return nil, fmt.Errorf("the parameter key '%s' not found", part)
		}
		updatedSubVal, err := updateNestedStructure(subVal, parts[1:], value)
		if err != nil {
			return nil, err
		}
		cv[part] = updatedSubVal
		return cv, nil
	default:
		return nil, fmt.Errorf("unsupported type at '%s'", part)
	}
}
