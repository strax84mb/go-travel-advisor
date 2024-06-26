// Code generated by ffjson <https://github.com/pquerna/ffjson>. DO NOT EDIT.
// source: save_comment.go

package dto

import (
	"bytes"
	"fmt"
	fflib "github.com/pquerna/ffjson/fflib/v1"
)

const (
	ffjtSaveCommentDtobase = iota
	ffjtSaveCommentDtonosuchkey

	ffjtSaveCommentDtoCityID

	ffjtSaveCommentDtoPosterID

	ffjtSaveCommentDtoText
)

var ffjKeySaveCommentDtoCityID = []byte("cityId")

var ffjKeySaveCommentDtoPosterID = []byte("posterId")

var ffjKeySaveCommentDtoText = []byte("text")

// UnmarshalJSON umarshall json - template of ffjson
func (j *SaveCommentDto) UnmarshalJSON(input []byte) error {
	fs := fflib.NewFFLexer(input)
	return j.UnmarshalJSONFFLexer(fs, fflib.FFParse_map_start)
}

// UnmarshalJSONFFLexer fast json unmarshall - template ffjson
func (j *SaveCommentDto) UnmarshalJSONFFLexer(fs *fflib.FFLexer, state fflib.FFParseState) error {
	var err error
	currentKey := ffjtSaveCommentDtobase
	_ = currentKey
	tok := fflib.FFTok_init
	wantedTok := fflib.FFTok_init

mainparse:
	for {
		tok = fs.Scan()
		//	println(fmt.Sprintf("debug: tok: %v  state: %v", tok, state))
		if tok == fflib.FFTok_error {
			goto tokerror
		}

		switch state {

		case fflib.FFParse_map_start:
			if tok != fflib.FFTok_left_bracket {
				wantedTok = fflib.FFTok_left_bracket
				goto wrongtokenerror
			}
			state = fflib.FFParse_want_key
			continue

		case fflib.FFParse_after_value:
			if tok == fflib.FFTok_comma {
				state = fflib.FFParse_want_key
			} else if tok == fflib.FFTok_right_bracket {
				goto done
			} else {
				wantedTok = fflib.FFTok_comma
				goto wrongtokenerror
			}

		case fflib.FFParse_want_key:
			// json {} ended. goto exit. woo.
			if tok == fflib.FFTok_right_bracket {
				goto done
			}
			if tok != fflib.FFTok_string {
				wantedTok = fflib.FFTok_string
				goto wrongtokenerror
			}

			kn := fs.Output.Bytes()
			if len(kn) <= 0 {
				// "" case. hrm.
				currentKey = ffjtSaveCommentDtonosuchkey
				state = fflib.FFParse_want_colon
				goto mainparse
			} else {
				switch kn[0] {

				case 'c':

					if bytes.Equal(ffjKeySaveCommentDtoCityID, kn) {
						currentKey = ffjtSaveCommentDtoCityID
						state = fflib.FFParse_want_colon
						goto mainparse
					}

				case 'p':

					if bytes.Equal(ffjKeySaveCommentDtoPosterID, kn) {
						currentKey = ffjtSaveCommentDtoPosterID
						state = fflib.FFParse_want_colon
						goto mainparse
					}

				case 't':

					if bytes.Equal(ffjKeySaveCommentDtoText, kn) {
						currentKey = ffjtSaveCommentDtoText
						state = fflib.FFParse_want_colon
						goto mainparse
					}

				}

				if fflib.SimpleLetterEqualFold(ffjKeySaveCommentDtoText, kn) {
					currentKey = ffjtSaveCommentDtoText
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				if fflib.EqualFoldRight(ffjKeySaveCommentDtoPosterID, kn) {
					currentKey = ffjtSaveCommentDtoPosterID
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				if fflib.SimpleLetterEqualFold(ffjKeySaveCommentDtoCityID, kn) {
					currentKey = ffjtSaveCommentDtoCityID
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				currentKey = ffjtSaveCommentDtonosuchkey
				state = fflib.FFParse_want_colon
				goto mainparse
			}

		case fflib.FFParse_want_colon:
			if tok != fflib.FFTok_colon {
				wantedTok = fflib.FFTok_colon
				goto wrongtokenerror
			}
			state = fflib.FFParse_want_value
			continue
		case fflib.FFParse_want_value:

			if tok == fflib.FFTok_left_brace || tok == fflib.FFTok_left_bracket || tok == fflib.FFTok_integer || tok == fflib.FFTok_double || tok == fflib.FFTok_string || tok == fflib.FFTok_bool || tok == fflib.FFTok_null {
				switch currentKey {

				case ffjtSaveCommentDtoCityID:
					goto handle_CityID

				case ffjtSaveCommentDtoPosterID:
					goto handle_PosterID

				case ffjtSaveCommentDtoText:
					goto handle_Text

				case ffjtSaveCommentDtonosuchkey:
					err = fs.SkipField(tok)
					if err != nil {
						return fs.WrapErr(err)
					}
					state = fflib.FFParse_after_value
					goto mainparse
				}
			} else {
				goto wantedvalue
			}
		}
	}

handle_CityID:

	/* handler: j.CityID type=int64 kind=int64 quoted=false*/

	{
		if tok != fflib.FFTok_integer && tok != fflib.FFTok_null {
			return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for int64", tok))
		}
	}

	{

		if tok == fflib.FFTok_null {

		} else {

			tval, err := fflib.ParseInt(fs.Output.Bytes(), 10, 64)

			if err != nil {
				return fs.WrapErr(err)
			}

			j.CityID = int64(tval)

		}
	}

	state = fflib.FFParse_after_value
	goto mainparse

handle_PosterID:

	/* handler: j.PosterID type=int64 kind=int64 quoted=false*/

	{
		if tok != fflib.FFTok_integer && tok != fflib.FFTok_null {
			return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for int64", tok))
		}
	}

	{

		if tok == fflib.FFTok_null {

		} else {

			tval, err := fflib.ParseInt(fs.Output.Bytes(), 10, 64)

			if err != nil {
				return fs.WrapErr(err)
			}

			j.PosterID = int64(tval)

		}
	}

	state = fflib.FFParse_after_value
	goto mainparse

handle_Text:

	/* handler: j.Text type=string kind=string quoted=false*/

	{

		{
			if tok != fflib.FFTok_string && tok != fflib.FFTok_null {
				return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for string", tok))
			}
		}

		if tok == fflib.FFTok_null {

		} else {

			outBuf := fs.Output.Bytes()

			j.Text = string(string(outBuf))

		}
	}

	state = fflib.FFParse_after_value
	goto mainparse

wantedvalue:
	return fs.WrapErr(fmt.Errorf("wanted value token, but got token: %v", tok))
wrongtokenerror:
	return fs.WrapErr(fmt.Errorf("ffjson: wanted token: %v, but got token: %v output=%s", wantedTok, tok, fs.Output.String()))
tokerror:
	if fs.BigError != nil {
		return fs.WrapErr(fs.BigError)
	}
	err = fs.Error.ToError()
	if err != nil {
		return fs.WrapErr(err)
	}
	panic("ffjson-generated: unreachable, please report bug.")
done:

	return nil
}
