package headerenc

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/Deimvis/go-ext/go1.25/ext"
	"github.com/Deimvis/go-ext/go1.25/xhttp"
	"github.com/Deimvis/go-ext/go1.25/xptr"
	"github.com/stretchr/testify/require"
)

func TestUnmarshal_Decode(t *testing.T) {
	type h = map[string][]string
	type vs = []string
	type keyfield[T any] struct {
		Key T `header:"key"`
	}
	tcs := []struct {
		title   string
		headers map[string][]string
		out     any
		expOut  any
		expMsg  UnaryPredicateReflect[string]
	}{
		{
			"string",
			h{"key": vs{"value"}},
			&struct {
				Key string `header:"key"`
			}{},
			&struct {
				Key string `header:"key"`
			}{Key: "value"},
			nil,
		},
		{
			"int",
			h{"key": vs{"42"}},
			&struct {
				Key int `header:"key"`
			}{},
			&struct {
				Key int `header:"key"`
			}{Key: 42},
			nil,
		},
		{
			"float",
			h{"key": vs{"42.123"}},
			&struct {
				Key float64 `header:"key"`
			}{},
			&struct {
				Key float64 `header:"key"`
			}{Key: 42.123},
			nil,
		},
		{
			"slice",
			h{"key": vs{"value1", "value2"}},
			&struct {
				Key []string `header:"key"`
			}{},
			&struct {
				Key []string `header:"key"`
			}{Key: []string{"value1", "value2"}},
			nil,
		},
		{
			"scalar/no-values",
			h{"key": vs{}},
			&struct {
				Key *string `header:"key"`
			}{},
			&struct {
				Key *string `header:"key"`
			}{Key: nil},
			nil,
		},
		{
			"scalar/multiple-values",
			h{"key": vs{"v1", "v2"}},
			&struct {
				Key *string `header:"key"`
			}{},
			nil,
			BindSecond(eqPred{}, "header key has multiple values, but expected single: key"),
		},
		{
			"container/no-values",
			h{"key": vs{}},
			&struct {
				Key []string `header:"key"`
			}{},
			&struct {
				Key []string `header:"key"`
			}{Key: nil},
			nil,
		},
		{
			"container/single-value",
			h{"key": vs{"value"}},
			&struct {
				Key []string `header:"key"`
			}{},
			&struct {
				Key []string `header:"key"`
			}{Key: []string{"value"}},
			nil,
		},
		{
			"container/multiple-value",
			h{"key": vs{"value1", "value2"}},
			&struct {
				Key []string `header:"key"`
			}{},
			&struct {
				Key []string `header:"key"`
			}{Key: []string{"value1", "value2"}},
			nil,
		},
		{
			"matched/exact",
			h{"key": vs{"value"}},
			&struct {
				Key string `header:"key"`
			}{},
			&struct {
				Key string `header:"key"`
			}{Key: "value"},
			nil,
		},
		{
			"matched/1-lettercase-diff",
			h{"Key": vs{"value"}},
			&struct {
				Key string `header:"key"`
			}{},
			&struct {
				Key string `header:"key"`
			}{Key: "value"},
			nil,
		},
		{
			"matched/many-lettercase-diff",
			h{"KeY": vs{"value"}},
			&struct {
				Key string `header:"key"`
			}{},
			&struct {
				Key string `header:"key"`
			}{Key: "value"},
			nil,
		},
		{
			"not-matched-headers-ignored",
			h{"key": vs{"value"}, "other-key": vs{"other-value"}},
			&struct {
				Key string `header:"key"`
			}{},
			&struct {
				Key string `header:"key"`
			}{Key: "value"},
			nil,
		},
		{
			"not-matched-struct-fields-ignored",
			h{"key": vs{"value"}},
			&struct {
				Key      string `header:"key"`
				OtherKey string `header:"other-key"`
			}{},
			&struct {
				Key      string `header:"key"`
				OtherKey string `header:"other-key"`
			}{Key: "value", OtherKey: ""},
			nil,
		},
		{
			"invalid-out/type-string",
			h{"key": vs{"value"}},
			xptr.T("out"),
			nil,
			BindSecond(eqPred{}, "unsupported input kind: string"),
		},
		{
			"invalid-out/header-name-collision-direct",
			h{"key": vs{"value"}},
			&struct {
				Key1 string `header:"key"`
				Key2 string `header:"key"`
			}{},
			nil,
			All(
				BindSecond(substrPred{}, "multiple fields in struct"),
				BindSecond(substrPred{}, "with same name key"),
			),
		},
		{
			"invalid-out/header-name-collision-direct-case-ins",
			h{"key": vs{"value"}},
			&struct {
				Key1 string `header:"key"`
				Key2 string `header:"Key"`
			}{},
			nil,
			All(
				BindSecond(substrPred{}, "multiple fields in struct"),
				BindSecond(substrPred{}, "with same name key"),
			),
		},
		{
			"invalid-out/header-name-collision-with-embed",
			h{"key": vs{"value"}},
			&struct {
				keyfield[string]
				Key2 string `header:"key"`
			}{},
			nil,
			All(
				BindSecond(substrPred{}, "multiple fields in struct"),
				BindSecond(substrPred{}, "with same name key"),
			),
		},
	}
	for _, tc := range tcs {
		t.Run(tc.title, func(t *testing.T) {
			headers := xhttp.AsConstHeader(tc.headers)
			err := Unmarshal(headers, tc.out, nil)
			if tc.expMsg != nil {
				require.Error(t, err)
				ok := tc.expMsg.Pred()(err.Error())
				require.True(t, ok, tc.expMsg.CallExplanationWithValues(ok, err.Error()))
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.out, tc.expOut)
			}
		})
	}
}

func TestUnmarshal_NotMatchedDecode(t *testing.T) {
	type h = map[string][]string
	type vs = []string
	type simple[T any] struct {
		Key T `header:"key"`
	}
	tcs := []struct {
		title            string
		headers          map[string][]string
		out              any
		expNotMatchedOut map[string][]string
		expMsg           UnaryPredicateReflect[string]
	}{
		{
			"all-matched",
			h{"key": vs{"value"}},
			&struct {
				Key string `header:"key"`
			}{},
			h{},
			nil,
		},
		{
			"some-not-matched",
			h{"key": vs{"value"}, "other-key": vs{"other-value"}},
			&struct {
				Key string `header:"key"`
			}{},
			h{"other-key": vs{"other-value"}},
			nil,
		},
	}
	for _, tc := range tcs {
		t.Run(tc.title, func(t *testing.T) {
			headers := xhttp.AsConstHeader(tc.headers)
			notMatchedOut := make(http.Header)
			err := Unmarshal(headers, tc.out, notMatchedOut)
			if tc.expMsg != nil {
				require.Error(t, err)
				ok := tc.expMsg.Pred()(err.Error())
				require.True(t, ok, tc.expMsg.CallExplanationWithValues(ok, err.Error()))
			} else {
				require.NoError(t, err)
				// normalize
				expNotMatchedOut := make(http.Header)
				for k, vs := range tc.expNotMatchedOut {
					for _, v := range vs {
						expNotMatchedOut.Add(k, v)
					}
				}
				for k, exp := range expNotMatchedOut {
					act := notMatchedOut.Values(k)
					require.True(t, len(act) > 0, "non-matched headers do not include header %s", k)
					require.Equal(t, exp, act, "non-matched headers have wrong value for %s", k)
				}
				for k, act := range notMatchedOut {
					exp := expNotMatchedOut.Values(k)
					require.True(t, len(exp) > 0, "non-matched headers include extra header %s", k)
					require.Equal(t, exp, act, "non-matched headers have wrong value for %s", k)
				}
			}
		})
	}
}

// TODO: use xcheck (xwhether) instead when predicates will be supported
// type Fact[T any] func() bool
type ZeroPredicate func() bool // aka Fact
type UnaryPredicate[T any] func(T) bool
type BinaryPredicate[T1, T2 any] func(T1, T2) bool

type UnaryPredicateCallState interface {
	// Call-specific explanation vs general explanation in predicateReflect.Explanation()
	// Shows exact reason why failed.
	Explanation() string
	ExplanationWithValues() string
}

type UnaryPredicateObservable[T any, State UnaryPredicateCallState] func(T) (bool, State)

// TODO:
// func (up UnaryPredicate[T]) ToReflect() UnaryPredicateReflect[T] {
// 	return nil
// }

type UnaryPredicateReflect[T any] interface {
	Pred() UnaryPredicate[T]
	// Description Negative? or Not(unaryPred) .Description() { Sprintf("not(%s)", subpred.Description()) } ? but what print on fail?
	Description() string
	DescriptionWithValues(T) string

	CallExplanation(bool) string
	CallExplanationWithValues(bool, T) string
}

type UnaryPredicateObservableReflect[T any, State UnaryPredicateCallState] interface {
	PredObservable() UnaryPredicateObservable[T, State]
}

type BinaryPredicateReflect[T1, T2 any] interface {
	Pred() BinaryPredicate[T1, T2]
	Description() string
	DescriptionWithValues(T1, T2) string

	CallExplanation(bool) string
	CallExplanationWithValues(bool, T1, T2) string
}

// xbpred (binary predicate)
// func All(bpos ...BinaryPredicateReflect[T1, T2 any])

// xupred (unary predicate)
func All[T any](upos ...UnaryPredicateReflect[T]) UnaryPredicateReflect[T] {
	return allUnary[T]{upos: upos}
}

type allUnary[T any] struct {
	upos []UnaryPredicateReflect[T]
}

func (au allUnary[T]) Pred() UnaryPredicate[T] {
	return func(v T) bool {
		for _, upo := range au.upos {
			if !upo.Pred()(v) {
				return false
			}
		}
		return true
	}
}

func (au allUnary[T]) Description() string {
	return fmt.Sprintf("all(%s)", strings.Join(
		ext.Map(au.upos, func(upo UnaryPredicateReflect[T]) string {
			return upo.Description()
		}),
		",",
	))
}

func (au allUnary[T]) DescriptionWithValues(v T) string {
	// TODO: Pred() returns Reflect predicate with method DescriptionWithValues() that doesn't require values
	// so
	return fmt.Sprintf("some of checks(%s) is not true", strings.Join(
		ext.Map(au.upos, func(upo UnaryPredicateReflect[T]) string {
			return upo.DescriptionWithValues(v)
		}),
		",",
	))
}

func (au allUnary[T]) CallExplanation(out bool) string {
	if out {
		return fmt.Sprintf("all(%s)", strings.Join(
			ext.Map(au.upos, func(upo UnaryPredicateReflect[T]) string {
				return upo.CallExplanation(out)
			}),
			",",
		))
	}
	return fmt.Sprintf("not all(%s)", strings.Join(
		ext.Map(au.upos, func(upo UnaryPredicateReflect[T]) string {
			return upo.CallExplanation(out)
		}),
		",",
	))
}

func (au allUnary[T]) CallExplanationWithValues(out bool, v T) string {
	if out {
		return fmt.Sprintf("all(%s)", strings.Join(
			ext.Map(au.upos, func(upo UnaryPredicateReflect[T]) string {
				return upo.CallExplanationWithValues(out, v)
			}),
			",",
		))
	}
	return fmt.Sprintf("not all(%s)", strings.Join(
		ext.Map(au.upos, func(upo UnaryPredicateReflect[T]) string {
			return upo.CallExplanationWithValues(out, v)
		}),
		",",
	))
}

// func All(binaryPreds...) BinaryPred
// func All(unaryPreds...) UnaryPred

type eqPred struct{}

func (p eqPred) Pred() BinaryPredicate[string, string] {
	return func(v1 string, v2 string) bool {
		return v1 == v2
	}
}

func (p eqPred) Description() string {
	return "equal"
}

func (p eqPred) DescriptionWithValues(v1 string, v2 string) string {
	return fmt.Sprintf("'%s' != '%s'", v1, v2)
}

func (p eqPred) CallExplanation(out bool) string {
	if out {
		return "equal"
	}
	return "not equal"
}

func (p eqPred) CallExplanationWithValues(out bool, v1 string, v2 string) string {
	if out {
		return fmt.Sprintf("'%s' == '%s'", v1, v2)
	}
	return fmt.Sprintf("'%s' != '%s'", v1, v2)
}

type substrPred struct{}

func (p substrPred) Pred() BinaryPredicate[string, string] {
	return func(v1 string, v2 string) bool {
		return strings.Contains(v1, v2)
	}
}

func (p substrPred) Description() string {
	return "has substring"
}

func (p substrPred) DescriptionWithValues(v1 string, v2 string) string {
	return fmt.Sprintf("'%s' has substring '%s'", v1, v2)
}

func (p substrPred) CallExplanation(out bool) string {
	if out {
		return "has substring"
	}
	return "has no substring"
}

func (p substrPred) CallExplanationWithValues(out bool, v1 string, v2 string) string {
	if out {
		return fmt.Sprintf("'%s' has substring '%s'", v1, v2)
	}
	return fmt.Sprintf("'%s' has no substring '%s'", v1, v2)
}

func BindFirst[T1, T2 any](bpo BinaryPredicateReflect[T1, T2], v1 T1) UnaryPredicateReflect[T2] {
	return unaryPredInplace[T2]{
		up: func(v T2) bool {
			return bpo.Pred()(v1, v)
		},
		desc: func() string {
			return fmt.Sprintf("(%s).bind(<placeholder>, _)", bpo.Description())
		},
		descWithValues: func(v T2) string {
			return bpo.DescriptionWithValues(v1, v)
		},
		expl: func(out bool) string {
			return bpo.CallExplanation(out)
		},
		explWithValues: func(out bool, v T2) string {
			return bpo.CallExplanationWithValues(out, v1, v)
		},
	}
}

func BindSecond[T1, T2 any](bpo BinaryPredicateReflect[T1, T2], v2 T2) UnaryPredicateReflect[T1] {
	return unaryPredInplace[T1]{
		up: func(v T1) bool {
			return bpo.Pred()(v, v2)
		},
		desc: func() string {
			return fmt.Sprintf("(%s).bind(_, <placeholder>)", bpo.Description())
		},
		descWithValues: func(v T1) string {
			return bpo.DescriptionWithValues(v, v2)
		},
		expl: func(out bool) string {
			return bpo.CallExplanation(out)
		},
		explWithValues: func(out bool, v T1) string {
			return bpo.CallExplanationWithValues(out, v, v2)
		},
	}
}

type unaryPredInplace[T any] struct {
	up             UnaryPredicate[T]
	desc           func() string
	descWithValues func(T) string
	expl           func(bool) string
	explWithValues func(bool, T) string
}

func (upi unaryPredInplace[T]) Pred() UnaryPredicate[T] {
	return upi.up
}

func (upi unaryPredInplace[T]) Description() string {
	return upi.desc()
}

func (upi unaryPredInplace[T]) DescriptionWithValues(v T) string {
	return upi.descWithValues(v)
}

func (upi unaryPredInplace[T]) CallExplanation(out bool) string {
	return upi.expl(out)
}

func (upi unaryPredInplace[T]) CallExplanationWithValues(out bool, v T) string {
	return upi.explWithValues(out, v)
}
