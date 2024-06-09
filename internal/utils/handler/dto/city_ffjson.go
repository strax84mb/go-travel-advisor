// Code generated by ffjson <https://github.com/pquerna/ffjson>. DO NOT EDIT.
// source: city.go

package dto

import (
	fflib "github.com/pquerna/ffjson/fflib/v1"
)

// MarshalJSON marshal bytes to json - template
func (j *CityAirportDto) MarshalJSON() ([]byte, error) {
	var buf fflib.Buffer
	if j == nil {
		buf.WriteString("null")
		return buf.Bytes(), nil
	}
	err := j.MarshalJSONBuf(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// MarshalJSONBuf marshal buff to json - template
func (j *CityAirportDto) MarshalJSONBuf(buf fflib.EncodingBuffer) error {
	if j == nil {
		buf.WriteString("null")
		return nil
	}
	var err error
	var obj []byte
	_ = obj
	_ = err
	buf.WriteString(`{"id":`)
	fflib.FormatBits2(buf, uint64(j.ID), 10, j.ID < 0)
	buf.WriteString(`,"name":`)
	fflib.WriteJsonString(buf, string(j.Name))
	buf.WriteByte('}')
	return nil
}

// MarshalJSON marshal bytes to json - template
func (j *CityDto) MarshalJSON() ([]byte, error) {
	var buf fflib.Buffer
	if j == nil {
		buf.WriteString("null")
		return buf.Bytes(), nil
	}
	err := j.MarshalJSONBuf(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// MarshalJSONBuf marshal buff to json - template
func (j *CityDto) MarshalJSONBuf(buf fflib.EncodingBuffer) error {
	if j == nil {
		buf.WriteString("null")
		return nil
	}
	var err error
	var obj []byte
	_ = obj
	_ = err
	buf.WriteString(`{ "id":`)
	fflib.FormatBits2(buf, uint64(j.ID), 10, j.ID < 0)
	buf.WriteString(`,"name":`)
	fflib.WriteJsonString(buf, string(j.Name))
	buf.WriteByte(',')
	if len(j.Airports) != 0 {
		buf.WriteString(`"airports":`)
		if j.Airports != nil {
			buf.WriteString(`[`)
			for i, v := range j.Airports {
				if i != 0 {
					buf.WriteString(`,`)
				}

				{

					err = v.MarshalJSONBuf(buf)
					if err != nil {
						return err
					}

				}
			}
			buf.WriteString(`]`)
		} else {
			buf.WriteString(`null`)
		}
		buf.WriteByte(',')
	}
	buf.Rewind(1)
	buf.WriteByte('}')
	return nil
}

// MarshalJSON marshal bytes to json - template
func (j *CityDtos) MarshalJSON() ([]byte, error) {
	var buf fflib.Buffer
	if j == nil {
		buf.WriteString("null")
		return buf.Bytes(), nil
	}
	err := j.MarshalJSONBuf(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// MarshalJSONBuf marshal buff to json - template
func (j *CityDtos) MarshalJSONBuf(buf fflib.EncodingBuffer) error {
	if j == nil {
		buf.WriteString("null")
		return nil
	}
	var err error
	var obj []byte
	_ = obj
	_ = err
	buf.WriteString(`{"items":`)
	if j.Items != nil {
		buf.WriteString(`[`)
		for i, v := range j.Items {
			if i != 0 {
				buf.WriteString(`,`)
			}

			{

				err = v.MarshalJSONBuf(buf)
				if err != nil {
					return err
				}

			}
		}
		buf.WriteString(`]`)
	} else {
		buf.WriteString(`null`)
	}
	buf.WriteByte('}')
	return nil
}