package vcjwt_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/alecthomas/assert/v2"
	"github.com/tbd54566975/web5-go/dids/didjwk"
	"github.com/tbd54566975/web5-go/jwt"
	"github.com/tbd54566975/web5-go/vc"
	"github.com/tbd54566975/web5-go/vc/vcjwt"
)

type vector struct {
	description string
	input       string
	errors      bool
}

func TestDecode(t *testing.T) {
	// TODO: move these to web5-spec repo test-vectors (Moe - 2024-02-24)
	vectors := []vector{
		{
			description: "fail to decode jwt",
			input:       "doodoo",
			errors:      true,
		},
		{
			description: "no claims",
			input:       "eyJhbGciOiJFZERTQSIsImtpZCI6ImRpZDpqd2s6ZXlKcmRIa2lPaUpQUzFBaUxDSmpjbllpT2lKRlpESTFOVEU1SWl3aWVDSTZJa3RPTVZGRU5ETkhZVkpJTm1OeWJpMTFTWFk0ZDBoT1pqZHlSMlkyVUY5RFMxZzJkbmsyTUdjMWQyc2lmUSMwIn0.e30.1iq9_pDtMlzL22h6xVY77nRNfXnR3oFU2kNYDAM52dPAs0l8zLL6AJ18B8rz9HziYzRo4Zo_jyYhq4nlHE3lBw",
			errors:      true,
		}, {
			description: "no vc claim",
			input:       "eyJhbGciOiJFZERTQSIsImtpZCI6ImRpZDpqd2s6ZXlKcmRIa2lPaUpQUzFBaUxDSmpjbllpT2lKRlpESTFOVEU1SWl3aWVDSTZJa0YyZFMxVVNsTlRWakZKT0dob1NrUmFlWGw1YUhnMmVHSlZjamRQT0RWMWRFMWpaa3RLVDJOblFWVWlmUSMwIn0.eyJoZWhlIjoiaGkifQ.QQ5aottVrsHRisxx7vRzin9CnyOcxeScxLOIy5qI30pV2FkXXBe3BdyujLS7i7M0CHW0eS9XhaVKe76504RZCQ",
			errors:      true,
		}, {
			description: "vc claim wrong type",
			input:       "eyJhbGciOiJFZERTQSIsImtpZCI6ImRpZDpqd2s6ZXlKcmRIa2lPaUpQUzFBaUxDSmpjbllpT2lKRlpESTFOVEU1SWl3aWVDSTZJbEZvUkhsdFprdzFaQzB0WjFGQ1prOVNOMlZFYjBrelprTjJOVUV0Y3pBMGFYZHlZMGRsTkVWd1lsa2lmUSMwIn0.eyJ2YyI6ImhpIn0.O_-xPUAZhi9W3OD1pJn4wN5Q9nZKYXcmtJPhWuk6WxlOXMca2jNXyjYpEKCJ1vFWZ4OHfSifErPvClLsH8-MCQ",
			errors:      true,
		}, {
			description: "legit",
			input:       "eyJhbGciOiJFZERTQSIsImtpZCI6ImRpZDpqd2s6ZXlKcmRIa2lPaUpQUzFBaUxDSmpjbllpT2lKRlpESTFOVEU1SWl3aWVDSTZJbVJwVGkxRlEydGhaRTVEVlVKZlUxRkhTVFJtVUdOZlluVmZObmt3VWpKRFdEUllkMjlTUzBjNFVEZ2lmUSMwIn0.eyJpc3MiOiJkaWQ6andrOmV5SnJkSGtpT2lKUFMxQWlMQ0pqY25ZaU9pSkZaREkxTlRFNUlpd2llQ0k2SW1ScFRpMUZRMnRoWkU1RFZVSmZVMUZIU1RSbVVHTmZZblZmTm5rd1VqSkRXRFJZZDI5U1MwYzRVRGdpZlEiLCJqdGkiOiJ1cm46dmM6dXVpZDoxOGQ5OTZjZi03N2YwLTRkYjgtOGQ5MS0zNGI1ZDY1NzcwNmUiLCJuYmYiOjE3MDg3NTY3ODUsInN1YiI6ImRpZDpqd2s6ZXlKcmRIa2lPaUpQUzFBaUxDSmpjbllpT2lKRlpESTFOVEU1SWl3aWVDSTZJbVJwVGkxRlEydGhaRTVEVlVKZlUxRkhTVFJtVUdOZlluVmZObmt3VWpKRFdEUllkMjlTUzBjNFVEZ2lmUSIsInZjIjp7IkBjb250ZXh0IjpbImh0dHBzOi8vd3d3LnczLm9yZy8yMDE4L2NyZWRlbnRpYWxzL3YxIl0sInR5cGUiOlsiVmVyaWZpYWJsZUNyZWRlbnRpYWwiXSwiaXNzdWVyIjoiZGlkOmp3azpleUpyZEhraU9pSlBTMUFpTENKamNuWWlPaUpGWkRJMU5URTVJaXdpZUNJNkltUnBUaTFGUTJ0aFpFNURWVUpmVTFGSFNUUm1VR05mWW5WZk5ua3dVakpEV0RSWWQyOVNTMGM0VURnaWZRIiwiY3JlZGVudGlhbFN1YmplY3QiOnsiaWQiOiJkaWQ6andrOmV5SnJkSGtpT2lKUFMxQWlMQ0pqY25ZaU9pSkZaREkxTlRFNUlpd2llQ0k2SW1ScFRpMUZRMnRoWkU1RFZVSmZVMUZIU1RSbVVHTmZZblZmTm5rd1VqSkRXRFJZZDI5U1MwYzRVRGdpZlEifSwiaWQiOiJ1cm46dmM6dXVpZDoxOGQ5OTZjZi03N2YwLTRkYjgtOGQ5MS0zNGI1ZDY1NzcwNmUiLCJpc3N1YW5jZURhdGUiOiIyMDI0LTAyLTI0VDA2OjM5OjQ1WiJ9fQ.Y1-9dFop7bg_0jvgZMLyE3CPjnSXH9SGTHeA_jn5HosYbhST8y_pK7LcDeCYLSgDfiIOeVsvJFqOr3XT2J2cDA",
			errors:      false,
		},
	}

	for _, tt := range vectors {
		t.Run(tt.description, func(t *testing.T) {
			decoded, err := vcjwt.Decode[vc.Claims](tt.input)

			if tt.errors == true {
				assert.Error(t, err)
				assert.Equal(t, vcjwt.Decoded[vc.Claims]{}, decoded)
			} else {
				assert.NoError(t, err)
				assert.NotEqual(t, vcjwt.Decoded[vc.Claims]{}, decoded)
			}
		})
	}
}

func TestDecode_SetClaims(t *testing.T) {
	issuer, err := didjwk.Create()
	assert.NoError(t, err)

	subject, err := didjwk.Create()
	assert.NoError(t, err)

	subjectClaims := vc.Claims{
		"firstName": "Randy",
		"lastName":  "McRando",
	}

	issuanceDate := time.Now().UTC()

	// missing issuer
	jwtClaims := jwt.Claims{
		JTI:        "abcd123",
		Issuer:     issuer.URI,
		Subject:    subject.URI,
		NotBefore:  issuanceDate.Unix(),
		Expiration: issuanceDate.Add(time.Hour).Unix(),
		Misc: map[string]any{
			"vc": vc.DataModel[vc.Claims]{
				CredentialSubject: subjectClaims,
			},
		},
	}

	vcJWT, err := jwt.Sign(jwtClaims, issuer)
	assert.NoError(t, err)

	decoded, err := vcjwt.Decode[vc.Claims](vcJWT)
	assert.NoError(t, err)

	assert.Equal(t, jwtClaims.JTI, decoded.VC.ID)
	assert.Equal(t, jwtClaims.Issuer, decoded.VC.Issuer)
	assert.Equal(t, jwtClaims.Subject, decoded.VC.CredentialSubject.GetID())
	assert.Equal(t, issuanceDate.Format(time.RFC3339), decoded.VC.IssuanceDate)
	assert.NotZero(t, decoded.VC.ExpirationDate)
}

func TestVerify(t *testing.T) {
	vectors := []vector{
		{
			description: "no id",
			input:       "eyJhbGciOiJFZERTQSIsImtpZCI6ImRpZDpqd2s6ZXlKcmRIa2lPaUpQUzFBaUxDSmpjbllpT2lKRlpESTFOVEU1SWl3aWVDSTZJbXB5VjFkVE0yZGxXRmR4UWxSa2JtRnZXV042WDFsT1p6aEdOR1pqYmxOTVRFOWhhbmw2VVdZME5ITWlmUSMwIn0.eyJleHAiOjI2NTQ5ODU4MDQsImlzcyI6ImRpZDpqd2s6ZXlKcmRIa2lPaUpQUzFBaUxDSmpjbllpT2lKRlpESTFOVEU1SWl3aWVDSTZJbXB5VjFkVE0yZGxXRmR4UWxSa2JtRnZXV042WDFsT1p6aEdOR1pqYmxOTVRFOWhhbmw2VVdZME5ITWlmUSIsIm5iZiI6MTcwODkwNTgwNCwic3ViIjoiZGlkOmp3azpleUpyZEhraU9pSlBTMUFpTENKamNuWWlPaUpGWkRJMU5URTVJaXdpZUNJNklqaDNOV00wYm5NeFkyeFRlWEJQVTNKNFNIcG9WVXN0V1hrMlFYQlZaMEpwTFVVNGRrRXhPRzkwTjFVaWZRIiwidmMiOnsiQGNvbnRleHQiOlsiaHR0cHM6Ly93d3cudzMub3JnLzIwMTgvY3JlZGVudGlhbHMvdjEiXSwidHlwZSI6WyJWZXJpZmlhYmxlQ3JlZGVudGlhbCJdLCJpc3N1ZXIiOiIiLCJjcmVkZW50aWFsU3ViamVjdCI6eyJmaXJzdE5hbWUiOiJSYW5keSIsImxhc3ROYW1lIjoiTWNSYW5kbyJ9LCJpc3N1YW5jZURhdGUiOiIifX0.lqDVmpV4Z156rMIxaj3cfyFeL_uK9YNpxJt3wbtIhEWRw1uapLcHhLKqbA1UoBiYjnoske8g3widw6qRVviWBw",
			errors:      true,
		}, {
			description: "no issuer",
			input:       "eyJhbGciOiJFZERTQSIsImtpZCI6ImRpZDpqd2s6ZXlKcmRIa2lPaUpQUzFBaUxDSmpjbllpT2lKRlpESTFOVEU1SWl3aWVDSTZJbFpoYkdoUFVGbDNWMDh6YjE5VFZETTNPRmRYZFdFeE1GbFdTblJvWkZCNWEzWlpja3BvWVRKMU1WVWlmUSMwIn0.eyJleHAiOjI2NTQ5ODU4MjUsImp0aSI6ImFiY2QxMjMiLCJuYmYiOjE3MDg5MDU4MjUsInN1YiI6ImRpZDpqd2s6ZXlKcmRIa2lPaUpQUzFBaUxDSmpjbllpT2lKRlpESTFOVEU1SWl3aWVDSTZJbFpoUjNCcE5tUjJXazh0TUY5MGFWcEdhM0JvZEhOcFEzUmtRVkY0WnpSSWRrZHJSV3B0Y21KS1kwMGlmUSIsInZjIjp7IkBjb250ZXh0IjpbImh0dHBzOi8vd3d3LnczLm9yZy8yMDE4L2NyZWRlbnRpYWxzL3YxIl0sInR5cGUiOlsiVmVyaWZpYWJsZUNyZWRlbnRpYWwiXSwiaXNzdWVyIjoiIiwiY3JlZGVudGlhbFN1YmplY3QiOnsiZmlyc3ROYW1lIjoiUmFuZHkiLCJsYXN0TmFtZSI6Ik1jUmFuZG8ifSwiaXNzdWFuY2VEYXRlIjoiIn19.myGiiuvshaZCwZXwAVLiOtazCcHXBvAvZy6xCkf8SRzsCeJXXYJACpMeAt9Q2Aiml3u0fW9Y7tH6uI9RQna2Dw",
			errors:      true,
		}, {
			description: "no issuance date",
			input:       "eyJhbGciOiJFZERTQSIsImtpZCI6ImRpZDpqd2s6ZXlKcmRIa2lPaUpQUzFBaUxDSmpjbllpT2lKRlpESTFOVEU1SWl3aWVDSTZJbU0xTVcwMFdYWnJUekJtWWxGdU4ycFVTVzVTYUVvek9Xd3hRa2xWVVdWMFJuSnFVVlF3TjNkalJqUWlmUSMwIn0.eyJleHAiOjI2NTQ5ODU4NDMsImlzcyI6ImRpZDpqd2s6ZXlKcmRIa2lPaUpQUzFBaUxDSmpjbllpT2lKRlpESTFOVEU1SWl3aWVDSTZJbU0xTVcwMFdYWnJUekJtWWxGdU4ycFVTVzVTYUVvek9Xd3hRa2xWVVdWMFJuSnFVVlF3TjNkalJqUWlmUSIsImp0aSI6ImFiY2QxMjMiLCJzdWIiOiJkaWQ6andrOmV5SnJkSGtpT2lKUFMxQWlMQ0pqY25ZaU9pSkZaREkxTlRFNUlpd2llQ0k2SW5WWVJqY3lZMGhMVUhGWVVuTTVjbmRpYlRselpqSTRXRVpaUWtnM1NpMXJObmxVWlhoNmEwWXhibFVpZlEiLCJ2YyI6eyJAY29udGV4dCI6WyJodHRwczovL3d3dy53My5vcmcvMjAxOC9jcmVkZW50aWFscy92MSJdLCJ0eXBlIjpbIlZlcmlmaWFibGVDcmVkZW50aWFsIl0sImlzc3VlciI6IiIsImNyZWRlbnRpYWxTdWJqZWN0Ijp7ImZpcnN0TmFtZSI6IlJhbmR5IiwibGFzdE5hbWUiOiJNY1JhbmRvIn0sImlzc3VhbmNlRGF0ZSI6IiJ9fQ.hdJPrCyaI_9sMi7fsENTTI2I9UW6zOswHsao82baZFx0tcI8HRFf8PfeMQDIoInOzS09LHRaQNJ528ewzUYVCw",
			errors:      true,
		}, {
			description: "issuance date in future",
			input:       "eyJhbGciOiJFZERTQSIsImtpZCI6ImRpZDpqd2s6ZXlKcmRIa2lPaUpQUzFBaUxDSmpjbllpT2lKRlpESTFOVEU1SWl3aWVDSTZJbDlOZGtkRVNVTXpZMmd4VjFCaFZGRmlXbFZoZGpOS2QwUlRORk0zWVhaS1ZEZEZlRWsyVkVSblIwa2lmUSMwIn0.eyJleHAiOjQ1NDcxNDg1OTcsImlzcyI6ImRpZDpqd2s6ZXlKcmRIa2lPaUpQUzFBaUxDSmpjbllpT2lKRlpESTFOVEU1SWl3aWVDSTZJbDlOZGtkRVNVTXpZMmd4VjFCaFZGRmlXbFZoZGpOS2QwUlRORk0zWVhaS1ZEZEZlRWsyVkVSblIwa2lmUSIsImp0aSI6ImFiY2QxMjMiLCJuYmYiOjI2NTQ5ODg1OTcsInN1YiI6ImRpZDpqd2s6ZXlKcmRIa2lPaUpQUzFBaUxDSmpjbllpT2lKRlpESTFOVEU1SWl3aWVDSTZJa3BWUzNwaFFWaFVVVkZoVFZjNVVVSlpZVkpKTFZKSU5tZzNhRU5FTTFCaVFWWmlNMDVXTTAxc2NXY2lmUSIsInZjIjp7IkBjb250ZXh0IjpbImh0dHBzOi8vd3d3LnczLm9yZy8yMDE4L2NyZWRlbnRpYWxzL3YxIl0sInR5cGUiOlsiVmVyaWZpYWJsZUNyZWRlbnRpYWwiXSwiaXNzdWVyIjoiIiwiY3JlZGVudGlhbFN1YmplY3QiOnsiZmlyc3ROYW1lIjoiUmFuZHkiLCJsYXN0TmFtZSI6Ik1jUmFuZG8ifSwiaXNzdWFuY2VEYXRlIjoiIn19.jj-VhKaYZQ3_6MtjHi-OLsXGQCuXKOCZ_mOGCsDg2EMedRIF8RdxfG0oZjJpeeujpndkj80WLFbQgm2F699mDQ",
			errors:      true,
		}, {
			description: "no context",
			input:       "eyJhbGciOiJFZERTQSIsImtpZCI6ImRpZDpqd2s6ZXlKcmRIa2lPaUpQUzFBaUxDSmpjbllpT2lKRlpESTFOVEU1SWl3aWVDSTZJbE5oY1VsMFFVcG9WblV3UkVWRVYybHNRbUptYlRNd1RtRk1WMWRqTWkxc2EyeElZMVpQT1d3NVNqUWlmUSMwIn0.eyJleHAiOjI2NTQ5ODU4NzMsImlzcyI6ImRpZDpqd2s6ZXlKcmRIa2lPaUpQUzFBaUxDSmpjbllpT2lKRlpESTFOVEU1SWl3aWVDSTZJbE5oY1VsMFFVcG9WblV3UkVWRVYybHNRbUptYlRNd1RtRk1WMWRqTWkxc2EyeElZMVpQT1d3NVNqUWlmUSIsImp0aSI6ImFiY2QxMjMiLCJuYmYiOjE3MDg5MDU4NzMsInN1YiI6ImRpZDpqd2s6ZXlKcmRIa2lPaUpQUzFBaUxDSmpjbllpT2lKRlpESTFOVEU1SWl3aWVDSTZJbUp5WWpkd2JYQnZWemN6YmpWZk1VVmFUbUozTjNoMlZVSnhOREpKVDFCck4zRmFhSGgzVlVoclgzY2lmUSIsInZjIjp7IkBjb250ZXh0IjpudWxsLCJ0eXBlIjpbIlZlcmlmaWFibGVDcmVkZW50aWFsIl0sImlzc3VlciI6IiIsImNyZWRlbnRpYWxTdWJqZWN0Ijp7ImZpcnN0TmFtZSI6IlJhbmR5IiwibGFzdE5hbWUiOiJNY1JhbmRvIn0sImlzc3VhbmNlRGF0ZSI6IiJ9fQ.cVzTVS0XoRL0lWdOoNe_NBrAR_0Xr5d_9OD0LABYvQgcnuLoH2MCa4Yh3p7R5-SMkchcIcKEMIqXhGzPxWD9CA",
			errors:      true,
		}, {
			description: "missing base context",
			input:       "eyJhbGciOiJFZERTQSIsImtpZCI6ImRpZDpqd2s6ZXlKcmRIa2lPaUpQUzFBaUxDSmpjbllpT2lKRlpESTFOVEU1SWl3aWVDSTZJbk5FUm5WTVpFMVdhMUpHY0VVNU5uQXpUbkp2U21kM01tMXZTMjVxWWw5bU5rdEtORzVWZDNWVGRrVWlmUSMwIn0.eyJleHAiOjI2NTQ5ODYwMjIsImlzcyI6ImRpZDpqd2s6ZXlKcmRIa2lPaUpQUzFBaUxDSmpjbllpT2lKRlpESTFOVEU1SWl3aWVDSTZJbk5FUm5WTVpFMVdhMUpHY0VVNU5uQXpUbkp2U21kM01tMXZTMjVxWWw5bU5rdEtORzVWZDNWVGRrVWlmUSIsImp0aSI6ImFiY2QxMjMiLCJuYmYiOjE3MDg5MDYwMjIsInN1YiI6ImRpZDpqd2s6ZXlKcmRIa2lPaUpQUzFBaUxDSmpjbllpT2lKRlpESTFOVEU1SWl3aWVDSTZJblZXUVdzMVJFZHpXakUwUXpjdGN6SmxWQzFKWWtZek9XOUVaMFZ4WVRKcmVUWnFUekZxU21VeFVWRWlmUSIsInZjIjp7IkBjb250ZXh0IjpbImh0dHBzOi8vc29tZWNvbnRleHQuY29tIl0sInR5cGUiOlsiVmVyaWZpYWJsZUNyZWRlbnRpYWwiXSwiaXNzdWVyIjoiIiwiY3JlZGVudGlhbFN1YmplY3QiOnsiZmlyc3ROYW1lIjoiUmFuZHkiLCJsYXN0TmFtZSI6Ik1jUmFuZG8ifSwiaXNzdWFuY2VEYXRlIjoiIn19.unSfKA5sZliFVHH7UrdiPvtMoExMGj465CBcHazICe6h2irQT7WNo2I4BsfUj0S3QmYQHKPwI9zQYbaNUBT_DQ",
			errors:      true,
		}, {
			description: "no type",
			input:       "eyJhbGciOiJFZERTQSIsImtpZCI6ImRpZDpqd2s6ZXlKcmRIa2lPaUpQUzFBaUxDSmpjbllpT2lKRlpESTFOVEU1SWl3aWVDSTZJbmRhTlU4NVpGRnNlVWhCUjBSUk1VOHhRVTlZY0ZOT04zWnVVRTlzV201MGFHMUVWRGwwYWxabmNUUWlmUSMwIn0.eyJleHAiOjI2NTQ5ODU5MTEsImlzcyI6ImRpZDpqd2s6ZXlKcmRIa2lPaUpQUzFBaUxDSmpjbllpT2lKRlpESTFOVEU1SWl3aWVDSTZJbmRhTlU4NVpGRnNlVWhCUjBSUk1VOHhRVTlZY0ZOT04zWnVVRTlzV201MGFHMUVWRGwwYWxabmNUUWlmUSIsImp0aSI6ImFiY2QxMjMiLCJuYmYiOjE3MDg5MDU5MTEsInN1YiI6ImRpZDpqd2s6ZXlKcmRIa2lPaUpQUzFBaUxDSmpjbllpT2lKRlpESTFOVEU1SWl3aWVDSTZJa0o0ZWpsRFdXVmlXbFJmZWtkV2FUWkNWbUUxV1VGSUxVSTRRekZqTTJ0MmJuRktkVFpzUjJsM1h6UWlmUSIsInZjIjp7IkBjb250ZXh0IjpbImh0dHBzOi8vd3d3LnczLm9yZy8yMDE4L2NyZWRlbnRpYWxzL3YxIl0sInR5cGUiOm51bGwsImlzc3VlciI6IiIsImNyZWRlbnRpYWxTdWJqZWN0Ijp7ImZpcnN0TmFtZSI6IlJhbmR5IiwibGFzdE5hbWUiOiJNY1JhbmRvIn0sImlzc3VhbmNlRGF0ZSI6IiJ9fQ.8r5EQCWhxaerwTBjRBkzalqD7IKlBxfFkd9vcehbzI8Ndzhl13jCeZlrl7PGFGpgwFscyEffghI3SgO3HbgxCQ",
			errors:      true,
		}, {
			description: "missing base type",
			input:       "eyJhbGciOiJFZERTQSIsImtpZCI6ImRpZDpqd2s6ZXlKcmRIa2lPaUpQUzFBaUxDSmpjbllpT2lKRlpESTFOVEU1SWl3aWVDSTZJamRCY0cxTmN6Z3dhMHRHWW04emNsUkJPV2hvY1dSUU5rWlFlamhVUkVoeE1GUXdSRGRMUTFNNVJXTWlmUSMwIn0.eyJleHAiOjI2NTQ5ODYwNjIsImlzcyI6ImRpZDpqd2s6ZXlKcmRIa2lPaUpQUzFBaUxDSmpjbllpT2lKRlpESTFOVEU1SWl3aWVDSTZJamRCY0cxTmN6Z3dhMHRHWW04emNsUkJPV2hvY1dSUU5rWlFlamhVUkVoeE1GUXdSRGRMUTFNNVJXTWlmUSIsImp0aSI6ImFiY2QxMjMiLCJuYmYiOjE3MDg5MDYwNjIsInN1YiI6ImRpZDpqd2s6ZXlKcmRIa2lPaUpQUzFBaUxDSmpjbllpT2lKRlpESTFOVEU1SWl3aWVDSTZJbkp3TTJ0aWVGbFFPVXRRZFVWWE5tUTBRV3hGYldOeE0wNUNYMkpzYUZGVllWWnFlVEJmZFc0MlRXY2lmUSIsInZjIjp7IkBjb250ZXh0IjpbImh0dHBzOi8vd3d3LnczLm9yZy8yMDE4L2NyZWRlbnRpYWxzL3YxIl0sInR5cGUiOlsiS25vd25DdXN0b21lckNyZWRlbnRpYWwiXSwiaXNzdWVyIjoiIiwiY3JlZGVudGlhbFN1YmplY3QiOnsiZmlyc3ROYW1lIjoiUmFuZHkiLCJsYXN0TmFtZSI6Ik1jUmFuZG8ifSwiaXNzdWFuY2VEYXRlIjoiIn19.irkHerKN78NRWdUyjSrg1WURYuzHBm54WU5sfO1Lj-FySSxt2TIYddhd1BE8V7pE0xDF1zNqkUQv5S_VAZKtCg",
			errors:      true,
		}, {
			description: "expired",
			input:       "eyJhbGciOiJFZERTQSIsImtpZCI6ImRpZDpqd2s6ZXlKcmRIa2lPaUpQUzFBaUxDSmpjbllpT2lKRlpESTFOVEU1SWl3aWVDSTZJbTlQYTNsa2NYcGFhV3RJVFdGUFprUlRXVUp2TkhGRVlWZE9PRnAzYWxSRGQyRmtZMVZHZG1kZmVVVWlmUSMwIn0.eyJleHAiOjE3MDg5MDYxMzQsImlzcyI6ImRpZDpqd2s6ZXlKcmRIa2lPaUpQUzFBaUxDSmpjbllpT2lKRlpESTFOVEU1SWl3aWVDSTZJbTlQYTNsa2NYcGFhV3RJVFdGUFprUlRXVUp2TkhGRVlWZE9PRnAzYWxSRGQyRmtZMVZHZG1kZmVVVWlmUSIsImp0aSI6ImFiY2QxMjMiLCJuYmYiOjE3MDg5MDI1MzQsInN1YiI6ImRpZDpqd2s6ZXlKcmRIa2lPaUpQUzFBaUxDSmpjbllpT2lKRlpESTFOVEU1SWl3aWVDSTZJbGRXZURkcFEyVTFUbGM0WVVSb1ZWOVRMVE4xTjNRMVIxWXlSQzE0VDJSaWVtRndjbms1ZWxkZk1Xc2lmUSIsInZjIjp7IkBjb250ZXh0IjpbImh0dHBzOi8vd3d3LnczLm9yZy8yMDE4L2NyZWRlbnRpYWxzL3YxIl0sInR5cGUiOlsiVmVyaWZpYWJsZUNyZWRlbnRpYWwiXSwiaXNzdWVyIjoiIiwiY3JlZGVudGlhbFN1YmplY3QiOnsiZmlyc3ROYW1lIjoiUmFuZHkiLCJsYXN0TmFtZSI6Ik1jUmFuZG8ifSwiaXNzdWFuY2VEYXRlIjoiIn19.Ca8JkSQnVA_hEctO_Me9nTXfgLivgYPhOR-Z2diPzInNS3xPFnvZQc3AlmhHqjjeOP1q8BUezw0cbqsck4LFCw",
			errors:      true,
		}, {
			description: "legit",
			input:       "eyJhbGciOiJFZERTQSIsImtpZCI6ImRpZDpqd2s6ZXlKcmRIa2lPaUpQUzFBaUxDSmpjbllpT2lKRlpESTFOVEU1SWl3aWVDSTZJakk1Vm5Wc2REQmxVRlIyVDBsSVZETlFZVTV6Ym1keldFTTVkVlU1Y0UwMGJHRjFRamcxVFhRMVlsa2lmUSMwIn0.eyJpc3MiOiJkaWQ6andrOmV5SnJkSGtpT2lKUFMxQWlMQ0pqY25ZaU9pSkZaREkxTlRFNUlpd2llQ0k2SWpJNVZuVnNkREJsVUZSMlQwbElWRE5RWVU1emJtZHpXRU01ZFZVNWNFMDBiR0YxUWpnMVRYUTFZbGtpZlEiLCJqdGkiOiJ1cm46dmM6dXVpZDo1OTQ4ODc0MS0wOWVmLTRlY2MtYWY0ZC0yYmFjN2MwNDY0OTIiLCJuYmYiOjE3MDg5MDY2OTIsInN1YiI6ImRpZDpqd2s6ZXlKcmRIa2lPaUpQUzFBaUxDSmpjbllpT2lKRlpESTFOVEU1SWl3aWVDSTZJbXg0VWtZeU5FMTNjVmsyZVdaRU1uSmhRMTlCZEhaNFNreDJhMk5YUWtNeVYzazJkR2RuYkVGcmIwMGlmUSIsInZjIjp7IkBjb250ZXh0IjpbImh0dHBzOi8vd3d3LnczLm9yZy8yMDE4L2NyZWRlbnRpYWxzL3YxIl0sInR5cGUiOlsiVmVyaWZpYWJsZUNyZWRlbnRpYWwiXSwiaXNzdWVyIjoiZGlkOmp3azpleUpyZEhraU9pSlBTMUFpTENKamNuWWlPaUpGWkRJMU5URTVJaXdpZUNJNklqSTVWblZzZERCbFVGUjJUMGxJVkROUVlVNXpibWR6V0VNNWRWVTVjRTAwYkdGMVFqZzFUWFExWWxraWZRIiwiY3JlZGVudGlhbFN1YmplY3QiOnsiZmlyc3ROYW1lIjoiUmFuZHkiLCJpZCI6ImRpZDpqd2s6ZXlKcmRIa2lPaUpQUzFBaUxDSmpjbllpT2lKRlpESTFOVEU1SWl3aWVDSTZJbXg0VWtZeU5FMTNjVmsyZVdaRU1uSmhRMTlCZEhaNFNreDJhMk5YUWtNeVYzazJkR2RuYkVGcmIwMGlmUSIsImxhc3ROYW1lIjoiTWNSYW5kbyJ9LCJpZCI6InVybjp2Yzp1dWlkOjU5NDg4NzQxLTA5ZWYtNGVjYy1hZjRkLTJiYWM3YzA0NjQ5MiIsImlzc3VhbmNlRGF0ZSI6IjIwMjQtMDItMjZUMDA6MTg6MTJaIn19.k7SzVjXghX_SYCkK35R6ulaBtjPSFc5fTtA1n70OpM16Fp_oitBZ2KQfLJfHl4YWnv1Z4hMItptijkGtEhWyDA",
			errors:      false,
		},
	}

	for _, tt := range vectors {
		t.Run(tt.description, func(t *testing.T) {
			_, err := vcjwt.Verify[vc.Claims](tt.input)

			if tt.errors == true {
				fmt.Printf("true error: %v\n", err)
				assert.Error(t, err)
			} else {
				fmt.Printf("false error: %v\n", err)
				assert.NoError(t, err)
			}
		})
	}
}

func TestSign(t *testing.T) {
	issuer, err := didjwk.Create()
	assert.NoError(t, err)

	subject, err := didjwk.Create()
	assert.NoError(t, err)

	claims := vc.Claims{"id": subject.URI, "name": "Randy McRando"}
	cred := vc.Create(claims)

	vcJWT, err := vcjwt.Sign(cred, issuer)
	assert.NoError(t, err)
	assert.NotZero(t, vcJWT)

	// TODO: make test more reliable by not depending on another function in this package (Moe - 2024-02-25)
	decoded, err := vcjwt.Verify[vc.Claims](vcJWT)

	assert.NoError(t, err)
	assert.NotZero(t, decoded)
}
