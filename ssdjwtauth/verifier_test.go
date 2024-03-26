package ssdjwtauth

import (
	"crypto"
	"reflect"
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

var (
	validPEMKeys = map[string][]byte{
		"key1": []byte(`-----BEGIN RSA PUBLIC KEY-----
MIIBCgKCAQEAn0P9T5ortQ03fS2A0nggHoHxc2oHm6H1bxo16Iot8/9iHPKIn1oq
CzKgpTqudGthHR/rS0mQRy7NYK8hO2Bqg8S45Qigc08/y/6l8fZDh+aFWPeQz+NH
LxVSTlXBtwBZXdZSaPsh23sYxKBbXGNLuFjI9E8Bb1SCuQiwrPT9Y2ZbjNU8xPtq
si/M7YGLKTlis5QDUv6JfU+Lo47To6jLTiDWNWKm0pb7Qwm4qCuQL3Bunyar/NHV
f59b4lju2oMp86F5sITmvTJ8hE6Tyq/N1T9zqy6X6A4Fl0miPPd/sNPEb9cyjn6J
ODYIi0VDOWoJ92NCur2LXL02FjcryDu+MQIDAQAB
-----END RSA PUBLIC KEY-----`),
		"key2": []byte(`-----BEGIN RSA PUBLIC KEY-----
MIIBCgKCAQEAsDYp2PGTczqY5AHOUK5VlklKSsy6kTILZMNRW1R4mPzryMl5lUJb
kzHPGALJx1J+g98MnXvIydVy0ZSOEA/J2Eg2qW2C79oLtu5SahHKcHISWS8qzu1v
6pw3TkbxIQkT7GJ90ziFpFO+h1321aiJc8sTqOp+K3OaXRq2lh6kP0sqDsyhAnau
C8azrm2aO4d5MDpyhfBDABb7Z+YGSLzAD93WQW7QvHdrb2A9NIRBZz3MMRmARe98
pHcca4lcglqAIAqrLoBHHjgaYcoyMbCsJ1lSnn6t0p779iXfHpK9MjexwxzWHaGB
tyMfhBEv2rnuuBjHzwgmn1MxAmqMTRbwCwIDAQAB
-----END RSA PUBLIC KEY-----`),
	}

	junkKeys = map[string][]byte{
		"badkey": []byte(`foo`),
	}
)

func Test_SetKeys(t *testing.T) {
	type args struct {
		pemkeys map[string][]byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"valid keys",
			args{
				validPEMKeys,
			},
			false,
		},
		{
			"badkeys",
			args{
				junkKeys,
			},
			true,
		},
	}
	for _, tt := range tests {
		v := Verifier{}
		t.Run(tt.name, func(t *testing.T) {
			err := v.SetKeys(tt.args.pemkeys)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetKeys() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestVerifier_KeyFunc(t *testing.T) {
	type args struct {
		Token *jwt.Token
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			"key is found",
			args{
				&jwt.Token{
					Header: map[string]interface{}{
						"kid": "key1",
					},
				},
			},
			[]byte("key1 contents"),
			false,
		},
		{
			"no kid",
			args{
				&jwt.Token{},
			},
			nil,
			true,
		},
		{
			"key is not found",
			args{
				&jwt.Token{
					Header: map[string]interface{}{
						"kid": "key99",
					},
				},
			},
			nil,
			true,
		},
		{
			"kdi is not a string",
			args{
				&jwt.Token{
					Header: map[string]interface{}{
						"kid": 5,
					},
				},
			},
			nil,
			true,
		},
	}

	v := &Verifier{
		Keys: map[string]crypto.PublicKey{
			"key1": []byte("key1 contents"),
			"key2": []byte("key2 contents"),
		},
	}

	keyfunc := v.KeyFunc()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := keyfunc(tt.args.Token)
			if (err != nil) != tt.wantErr {
				t.Errorf("Did not expect error: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Verifier.KeyFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}