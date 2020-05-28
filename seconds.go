// Copyright 2020 Matt Ho
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
package epoch

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// Seconds represents seconds since unix epoch
type Seconds int64

// Add duration to an int64
func (s Seconds) Add(d time.Duration) Seconds {
	return s + Seconds(d/time.Second)
}

// Int64 represents an int64
func (s Seconds) Int64() int64 {
	return int64(s)
}

// MarshalJSON implements json.Marshaler
func (s Seconds) MarshalJSON() ([]byte, error) {
	return json.Marshal(int64(s))
}

// UnmarshalJSON implements json.Unmarshaler
func (s *Seconds) UnmarshalJSON(data []byte) error {
	var i int64
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}

	*s = Seconds(i)

	return nil
}

// TimeInLocation returns time in default location
func (s Seconds) Time() time.Time {
	return s.TimeInLocation(time.Local)
}

// TimeInLocation returns time in specified location
func (s Seconds) TimeInLocation(loc *time.Location) time.Time {
	return time.Unix(s.Int64(), 0).In(loc)
}

// MarshalDynamoDBAttributeValue implements dynamodbattribute.Marshaler
func (s Seconds) MarshalDynamoDBAttributeValue(av *dynamodb.AttributeValue) error {
	av.N = aws.String(strconv.FormatInt(s.Int64(), 10))
	return nil
}

// UnmarshalDynamoDBAttributeValue implements dynamodbattribute.Unmarshaler
func (s *Seconds) UnmarshalDynamoDBAttributeValue(av *dynamodb.AttributeValue) error {
	if av == nil || av.N == nil {
		return nil
	}

	v, err := strconv.ParseInt(*av.N, 10, 64)
	if err != nil {
		return err
	}

	*s = Seconds(v)

	return nil
}

// Now returns current time expressed as seconds
func Now() Seconds {
	return From(time.Now())
}

// From returns the epoch seconds from a given time
func From(t time.Time) Seconds {
	return Seconds(t.Unix())
}
