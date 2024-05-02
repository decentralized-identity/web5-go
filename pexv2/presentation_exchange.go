package pexv2

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/PaesslerAG/jsonpath"
	"github.com/tbd54566975/web5-go/vc"
	jsonschema "github.com/xeipuuv/gojsonschema"
)

func SelectCredentials(vcJwts []string, pd PresentationDefinition) ([]string, error) {

	result := make([]string, 0)
	for _, inputDescriptor := range pd.InputDescriptors {
		matchedVcJwts, _ := selectCredentialsPerId(vcJwts, inputDescriptor)
		if len(matchedVcJwts) == 0 {
			return []string{}, nil
		}
		result = append(result, matchedVcJwts...)

	}

	result = dedupeResult(result)
	return result, nil
}

func dedupeResult(input []string) []string {
	sort.Strings(input)
	var result []string

	for i, item := range input {
		if i == 0 || input[i-1] != item {
			result = append(result, item)
		}
	}
	return result
}

func selectCredentialsPerId(vcJwts []string, inputDescriptor InputDescriptor) ([]string, error) {
	answer := make([]string, 0)
	tokenizedField := make([]TokenPath, 0)
	schema := map[string]interface{}{
		"$schema":    "http://json-schema.org/draft-07/schema#",
		"type":       "object",
		"properties": map[string]interface{}{},
		"required":   []string{},
	}

	for _, field := range inputDescriptor.Constraints.Fields {
		token := inputDescriptor.generateRandomToken()
		tokenizedField = append(tokenizedField, TokenPath{Token: token, Paths: field.Path})

		properties, ok := schema["properties"].(map[string]interface{})
		if !ok {
			fmt.Printf("unable to assert 'properties' as map[string]interface{}")
		}

		if field.Filter != nil {
			properties[token] = field.Filter
		} else {
			// null is omitted
			anyType := map[string]interface{}{
				"type": []string{"string", "number", "boolean", "object", "array"},
			}
			properties[token] = anyType
		}
		if required, ok := schema["required"].([]string); ok {
			required = append(required, token)
			schema["required"] = required
		}

	}
	fmt.Printf("TokenPaths: %v\n", tokenizedField)
	schemaJSON, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling schema:", err)
	}
	fmt.Printf("Schema: %s\n", string(schemaJSON))
	// schema only has one property

	for _, vcJwt := range vcJwts {
		decoded, _ := vc.Decode[vc.Claims](vcJwt)
		vcJson := getVcJson(decoded)

		selectionCandidate := make(map[string]interface{})
		// { "asdfghj": "mybtcaddress", "qwert": "mydogeaddress"}
		for _, tokenPath := range tokenizedField {
			for _, path := range tokenPath.Paths {
				value, err := jsonpath.Get(path, vcJson)
				if err != nil {
					continue
				}

				fmt.Printf("putting token %s and paths %s with value %s in pathMatchedCandidates\n", tokenPath.Token, tokenPath.Paths, value)

				selectionCandidate[tokenPath.Token] = value
				break
			}
		}

		fmt.Printf("Selection Candidate: %v\n", selectionCandidate)
		// { "asdfghj": "mybtcaddress", "qwert": "mydogeaddress"}

		fmt.Println("Properties exist, validating with schema")
		validationResult, err := validateWithSchema(schema, selectionCandidate)
		if err != nil {
			fmt.Println("Error validating schema:", err)
		}

		if validationResult.Valid() {
			fmt.Printf("The vcJWT is valid against schema!!!!\n\n\n\n")
			answer = append(answer, vcJwt)
		}

	}

	return answer, nil

}

func validateWithSchema(schema map[string]interface{}, pathMatchedCandidates map[string]interface{}) (*jsonschema.Result, error) {
	schemaLoader := getSchemaLoader(schema, pathMatchedCandidates)
	documentLoader := jsonschema.NewGoLoader(pathMatchedCandidates)

	result, err := jsonschema.Validate(schemaLoader, documentLoader)
	return result, err
}

func getSchemaLoader(schema map[string]interface{}, selectionCandidates map[string]interface{}) jsonschema.JSONLoader {
	schemaJSON, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling schema:", err)
	}
	fmt.Printf("Schema: %s\n", string(schemaJSON))
	fmt.Printf("Selection Candidates: %v\n", selectionCandidates)

	schemaLoader := jsonschema.NewStringLoader(string(schemaJSON))
	return schemaLoader
}

func getVcJson(decoded vc.DecodedVCJWT[vc.Claims]) interface{} {
	marshaledVcJwt, _ := json.Marshal(decoded.JWT.Claims)
	var jsondata interface{}
	err := json.Unmarshal(marshaledVcJwt, &jsondata)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return interface{}(nil)
	}
	return jsondata
}
