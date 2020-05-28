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
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

func TestSeconds_JSON(t *testing.T) {
	t.Run("bad", func(t *testing.T) {
		var s Seconds
		err := (&s).UnmarshalJSON([]byte(`bad`))

		var want *json.SyntaxError
		ok := errors.As(err, &want)
		if !ok {
			t.Fatalf("got %T; want %T", err, want)
		}
	})

	t.Run("null", func(t *testing.T) {
		const text = `null`
		var s Seconds
		err := json.Unmarshal([]byte(text), &s)
		if err != nil {
			t.Fatalf("got %v; want nil", err)
		}
		if got, want := s, Seconds(0); got != want {
			t.Fatalf("got %v; want %v", got, want)
		}
	})

	t.Run("value", func(t *testing.T) {
		want := Seconds(123)
		data, err := json.Marshal(want)
		if err != nil {
			t.Fatalf("got %v; want nil", err)
		}

		var got Seconds
		err = json.Unmarshal(data, &got)
		if got != want {
			t.Fatalf("got %v; want %v", got, want)
		}
	})
}

func TestSeconds_Int64(t *testing.T) {
	s := Seconds(123)
	if got, want := s.Int64(), int64(123); got != want {
		t.Fatalf("got %v; want %v", got, want)
	}
}

func TestSeconds_Time(t *testing.T) {
	want := time.Now().In(time.Local).Round(time.Second)
	got := Seconds(want.Unix())
	assert.Equal(t, want, got.Time())
}

func TestSeconds_Add(t *testing.T) {
	var s Seconds = 123
	got := s.Add(time.Second)
	assert.Equal(t, s+1, got)
}

func TestSeconds_UnmarshalDynamoDBAttributeValue(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var s Seconds

		err := (&s).UnmarshalDynamoDBAttributeValue(nil)
		if err != nil {
			t.Fatalf("got %v; want nil", err)
		}

		err = (&s).UnmarshalDynamoDBAttributeValue(&dynamodb.AttributeValue{})
		if err != nil {
			t.Fatalf("got %v; want nil", err)
		}
	})

	t.Run("bad", func(t *testing.T) {
		var s Seconds

		err := (&s).UnmarshalDynamoDBAttributeValue(&dynamodb.AttributeValue{N: aws.String("bad")})

		var e *strconv.NumError
		ok := errors.As(err, &e)
		if !ok {
			t.Fatalf("got %T; want %T", err, e)
		}
	})

	t.Run("ok", func(t *testing.T) {
		want := Seconds(123)

		av, err := dynamodbattribute.Marshal(want)
		if err != nil {
			t.Fatalf("got %v; want nil", err)
		}

		var got Seconds
		err = (&got).UnmarshalDynamoDBAttributeValue(av)
		if err != nil {
			t.Fatalf("got %v; want nil", err)
		}
		if got != want {
			t.Fatalf("got %v; want %v", got, want)
		}
	})
}

func TestNow(t *testing.T) {
	s := Now()
	n := From(time.Now())
	got := n - s
	if got > 1 {
		t.Fatalf("got %v; want <= 1", got)
	}
}
