package fv

import (
	_ "embed"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed test/classic.html
var parsedTemplate []byte

func TestEntry(t *testing.T) {
	entry := Entry{
		Amount: 2,
		Price:  100,
		Vat:    8,
		Name:   "entry",
	}

	assert.Equal(t, float64(200), entry.GetNetPrice())
	assert.Equal(t, float64(16), entry.GetVatPrice())
	assert.Equal(t, float64(216), entry.GetGrossPrice())
}

func TestConfig(t *testing.T) {
	os.Setenv(FV_CONFIG_PATH_NAME, "test/test-config.yaml")

	t.Run("NewConfig", func(t *testing.T) {
		c, err := NewConfig()
		assert.Nil(t, err)
		assert.Equal(t, "Recipient Name", c.Recipient.Name)
		assert.Equal(t, "Recipient Address", c.Recipient.Address)
		assert.Equal(t, "123 456 78 90", c.Recipient.VatNumber)
		assert.Equal(t, "Seller Name", c.Seller.Name)
		assert.Equal(t, "Seller Address", c.Seller.Address)
		assert.Equal(t, "Seller City", c.Seller.City)
		assert.Equal(t, "098 765 43 21", c.Seller.VatNumber)
		assert.Equal(t, "12 1234 1234 0000 1234 1234 1234", c.Seller.BankNumber)
		assert.Equal(t, "Invoice Prefix", c.FV.HeaderPrefix)
		assert.Equal(t, "0000/1111/2222", c.FV.NO)
		assert.Equal(t, "2012-05-31", c.Now.Time().Format(dateFormat))
		assert.Equal(t, 2, len(c.FV.Entries))
		assert.Equal(t, "Entry [0] name", c.FV.Entries[0].Name)
		assert.Equal(t, 1000.12, c.FV.Entries[0].Price)
		assert.Equal(t, "szt", c.FV.Entries[0].Unit)
		assert.Equal(t, 1, c.FV.Entries[0].Amount)
		assert.Equal(t, 23, c.FV.Entries[0].Vat)
		assert.Equal(t, "Entry&nbsp;[1]&nbsp;name", c.FV.Entries[1].Name)
		assert.Equal(t, 2000.22, c.FV.Entries[1].Price)
		assert.Equal(t, "szt", c.FV.Entries[1].Unit)
		assert.Equal(t, 2, c.FV.Entries[1].Amount)
		assert.Equal(t, 8, c.FV.Entries[1].Vat)
	})

	t.Run("parseTemplate", func(t *testing.T) {
		t.Run("classic.html", func(t *testing.T) {
			c, err := NewConfig()
			assert.Nil(t, err)

			template, err := c.parseTemplate()
			assert.Nil(t, err)

			assert.Equal(t, parsedTemplate, template)
		})
	})

	t.Run("getFileName", func(t *testing.T) {
		c, err := NewConfig()
		assert.Nil(t, err)

		assert.Equal(t, "FV-2012-05.pdf", c.getFileName())
	})
}
