package lib_test

import (
	"testing"
	"unit-test-demo/unittest3/lib"
)

func TestCalculateTotals(t *testing.T) {
	tests := []struct {
		name        string
		items       []lib.Item
		discountPct int
		taxPct      int
		want        lib.Totals
	}{
		{
			name: "simple no discount no tax",
			items: []lib.Item{
				{Qty: 2, UnitPriceCents: 1500}, // $15.00 x 2 = $30.00
				{Qty: 1, UnitPriceCents: 2500}, // $25.00
			},
			discountPct: 0,
			taxPct:      0,
			want: lib.Totals{
				Subtotal: 5500, // $55.00
				Discount: 0,
				Tax:      0,
				Total:    5500,
			},
		},
		{
			name: "with 10% discount and 11% tax",
			items: []lib.Item{
				{Qty: 3, UnitPriceCents: 1999}, // $19.99 x 3 = 5997
				{Qty: 1, UnitPriceCents: 505},  // $5.05
			},
			discountPct: 10, // 10% of 6502 = 650.2 => 650 (halves-up; .2 rounds down)
			taxPct:      11, // 11% of (6502-650=5852) = 643.72 => 644 (halves-up)
			want: lib.Totals{
				Subtotal: 6502,
				Discount: 650,
				Tax:      644,
				Total:    6496, // 6502 - 650 + 644
			},
		},
		{
			name: "rounding halves-up edge",
			items: []lib.Item{
				{Qty: 1, UnitPriceCents: 333}, // $3.33
			},
			discountPct: 15, // 15% of 333 = 49.95 => 50 (halves-up)
			taxPct:      10, // 10% of (333-50=283) = 28.3 => 28
			want: lib.Totals{
				Subtotal: 333,
				Discount: 50,
				Tax:      28,
				Total:    311,
			},
		},
		// {
		// 	name: "invalid lines ignored",
		// 	items: []lib.Item{
		// 		{Qty: 0, UnitPriceCents: 1000},
		// 		{Qty: -1, UnitPriceCents: 2000},
		// 		{Qty: 2, UnitPriceCents: -500},
		// 		{Qty: 2, UnitPriceCents: 250}, // valid: 500
		// 	},
		// 	discountPct: 0,
		// 	taxPct:      0,
		// 	want: lib.Totals{
		// 		Subtotal: 500,
		// 		Discount: 0,
		// 		Tax:      0,
		// 		Total:    500,
		// 	},
		// },
	}

	for _, tt := range tests {
		tt := tt // capture
		t.Run(tt.name, func(t *testing.T) {
			got := lib.CalculateTotals(tt.items, tt.discountPct, tt.taxPct)
			if got != tt.want {
				t.Fatalf("got %+v, want %+v", got, tt.want)
			}

			// Purity sanity check: calling twice yields identical result.
			got2 := lib.CalculateTotals(tt.items, tt.discountPct, tt.taxPct)
			if got2 != got {
				t.Fatalf("non-deterministic result: first %+v, second %+v", got, got2)
			}
		})
	}
}

func TestRoundPercent(t *testing.T) {
	type args struct {
		amount lib.Money
		pct    int
	}
	cases := []struct {
		name string
		args args
		want lib.Money
	}{
		{"zero pct", args{amount: 1000, pct: 0}, 0},
		{"zero amount", args{amount: 0, pct: 5}, 0},
		{"round down", args{amount: 1000, pct: 1}, 10},  // 10.00
		{"round halves-up", args{amount: 333, pct: 15}, 50}, // 49.95 -> 50
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			if got := lib.RoundPercent(tc.args.amount, tc.args.pct); got != tc.want {
				t.Fatalf("roundPercent(%v,%v)=%v, want %v", tc.args.amount, tc.args.pct, got, tc.want)
			}
		})
	}
}
