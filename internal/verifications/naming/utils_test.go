package naming

import (
	"go/ast"
	"go/token"
	"reflect"
	"strings"
	"testing"

	"github.com/fdaines/arch-go/internal/model"

	"github.com/stretchr/testify/assert"
)

func TestNamingRuleUtils(t *testing.T) {
	t.Run("getPatternComparator case 1", func(t *testing.T) {
		comparator, s := getPatternComparator("")

		assert.Equal(t, reflect.ValueOf(comparator), reflect.ValueOf(strings.EqualFold))
		assert.Equal(t, "", s)
	})

	t.Run("getPatternComparator case 2", func(t *testing.T) {
		comparator, s := getPatternComparator("foobar*")

		assert.Equal(t, reflect.ValueOf(comparator), reflect.ValueOf(strings.HasPrefix))
		assert.Equal(t, "foobar", s)
	})

	t.Run("getPatternComparator case 3", func(t *testing.T) {
		comparator, s := getPatternComparator("*blablabla")

		assert.Equal(t, reflect.ValueOf(comparator), reflect.ValueOf(strings.HasSuffix))
		assert.Equal(t, "blablabla", s)
	})

	t.Run("getReturnValues case 1", func(t *testing.T) {
		resultsObject := &ast.FieldList{
			List: []*ast.Field{
				{
					Type: mockType{pos: 30, end: 40},
				},
				{
					Type: mockType{pos: 100, end: 130},
				},
			},
		}
		expectedValues := []string{"oigudsoigu", "vmcnvjkdfjgfdjgusfdjgh\nudfgjfd"}

		returnValues := getReturnValues(fileContent, resultsObject)

		assert.Equal(t, expectedValues, returnValues)
	})

	t.Run("getReturnValues case 2", func(t *testing.T) {
		resultsObject := &ast.FieldList{}
		returnValues := getReturnValues(fileContent, resultsObject)

		assert.Equal(t, []string{}, returnValues)
	})

	t.Run("getReturnValues case 3", func(t *testing.T) {
		returnValues := getReturnValues(fileContent, nil)

		assert.Equal(t, []string{}, returnValues)
	})

	t.Run("getParameters case 1", func(t *testing.T) {
		resultsObject := &ast.FieldList{
			List: []*ast.Field{
				{
					Type: mockType{pos: 50, end: 74},
				},
				{
					Type: mockType{pos: 200, end: 208},
				},
			},
		}
		expectedValues := []string{"suvicxxnvcxnvuceanckjdwn", "ewioufoi"}

		returnValues := getParameters(fileContent, resultsObject)

		assert.Equal(t, expectedValues, returnValues)
	})

	t.Run("getParameters case 2", func(t *testing.T) {
		resultsObject := &ast.FieldList{}
		returnValues := getParameters(fileContent, resultsObject)

		assert.Equal(t, []string{}, returnValues)
	})

	t.Run("getParameters case 3", func(t *testing.T) {
		returnValues := getParameters(fileContent, nil)

		assert.Equal(t, []string{}, returnValues)
	})

	t.Run("packageMustBeAnalyzed case 1", func(t *testing.T) {
		returnValue := packageMustBeAnalyzed(nil, "foo.**")

		assert.Equal(t, false, returnValue)
	})

	t.Run("packageMustBeAnalyzed case 2", func(t *testing.T) {
		pkg := &model.PackageInfo{Path: "foo/bar/blablabla/dummy"}
		returnValue := packageMustBeAnalyzed(pkg, "**.bar.**")

		assert.Equal(t, true, returnValue)
	})

	t.Run("packageMustBeAnalyzed case 3", func(t *testing.T) {
		pkg := &model.PackageInfo{Path: "foo/bar/blablabla/dummy"}
		returnValue := packageMustBeAnalyzed(pkg, "unknown.**")

		assert.Equal(t, false, returnValue)
	})

	t.Run("resolveStructName case 1", func(t *testing.T) {
		fd := &ast.FuncDecl{
			Recv: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.Ident{
							Name: "barfoo",
						},
					},
				},
			},
		}
		returnValue := resolveStructName(fd)

		assert.Equal(t, "barfoo", returnValue)
	})

	t.Run("resolveStructName case 2", func(t *testing.T) {
		fd := &ast.FuncDecl{
			Recv: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.UnaryExpr{},
					},
				},
			},
		}
		returnValue := resolveStructName(fd)

		assert.Equal(t, "", returnValue)
	})
}

type mockType struct {
	*ast.StructType
	pos token.Pos
	end token.Pos
}

func (m mockType) Pos() token.Pos { return m.pos }

func (m mockType) End() token.Pos { return m.end }

const fileContent = `
etuwoqtueowiqutietewioufoiduoigudsoiguoidsaguoidsuvicxxnvcxnvuceanckjdwncduwnosdncxnvzdjnvufewncjdvmcnvjkdfjgfdjgusfdjgh
udfgjfdshguseeirjfawkdopkodjsdijijfdjfioejvmdlksmvdivjiodsvetuwoqtueowiqutietewioufoiduoigudsoiguoidsaguoidsuvicxxnvcxnv
uceanckjdwnetuwoqtueowiqutietewioufoiduoigudsoiguoidsaguoidsuvicxxnvcxnvuceanckjdwnxxxnvcxnvuceanckjdwncduwnosdncxnvzdjn
vufewncjdvmcnvjkddsuvicxxnvcxnvuceanckjdwncduwnosdncxnvzdjnvufewncjfdvmcnvjkdfjgfdjgusfdjghudfgjfdshguseeirjfawkdopkodjs
dijijfdjfioejvmdlksmvdivjiodsvetuwoqtueowiqutietewioufoiduoigudsoiguoidsaguoidsuvicxxnvcxnvuceanckjdwnetuwoqtueowiqutiet
ewioufoiduoigudsoiguoidsaguoidsuvicxxnvcxnvuceanckjdwnxxxnvcxnvuceanckjdwncduwnosdncxnvzdjnvufewncjdvmcnvjkdfjgfdjgusfdj
ghudfgjfdshguseeirjfawkdopkodjsdijijfdjfioejvmdlksmvdivjiodfioejvmdlksmvdivjiodsvetuwoqtueowiqutietewioufoiduoigudsoiguo
uceanckjdwnetuwoqtueowiqutietewioufoiduoigudsoiguoidsaguoidsuvicxxnvcxnvuceanckjdwnxxxnvcxnvuceanckjdwncduwnosdncxnvzdjn
vufewncjdvmcnvjkddsuvicxxnvcxnvuceanckjdwncduwnosdncxnvzdjnvufewncjfdvmcnvjkdfjgfdjgusfdjghudfgjfdshguseeirjfawkdopkodjs
dijijfdjfioejvmdlksmvdivjiodsvetuwoqtueowiqutietewioufoiduoigudsoiguoidsaguoidsuvicxxnvcxnvuceanckjdwnetuwoqtueowiqutiet
ewioufoiduoigudsoiguoidsaguoidsuvicxxnvcxnvuceanckjdwnxxxnvcxnvuceanckjdwncduwnosdncxnvzdjnvufewncjdvmcnvjkdfjgfdjgusfdj
ghudfgjfdshguseeirjfawkdopkodjsdijijfdjfioejvmdlksmvdivjiodidsaguoidsuvicxxnvcxnvnvuceanckjdwncduwnosnvuceanckjdwncduwno
`

type mockStarExpr struct {
	*ast.StarExpr
	pos token.Pos
	end token.Pos
}

func (m mockStarExpr) Pos() token.Pos { return m.pos }

func (m mockStarExpr) End() token.Pos { return m.end }
