package lib

// Money is in cents to keep math exact (e.g., $12.34 => 1234).
type Money int64

// Item represents a purchasable line item.
type Item struct {
	Qty             int64 // e.g., 3
	UnitPriceCents  Money // e.g., 1999 for $19.99
}

// Totals contains the derived amounts from CalculateTotals.
type Totals struct {
	Subtotal Money // sum of qty * unit price
	Discount Money // rounded to nearest cent
	Tax      Money // rounded to nearest cent
	Total    Money // Subtotal - Discount + Tax
}

// CalculateTotals computes invoice totals using integer math.
// discountPct and taxPct are plain percentages (e.g., 10 => 10%).
// Rounding: halves up (0.5+) to the next cent.
func CalculateTotals(items []Item, discountPct, taxPct int) Totals {
	var subtotal Money
	for _, it := range items {
		if it.Qty <= 0 || it.UnitPriceCents < 0 {
			// Defensive: ignore nonsensical lines; still pure.
			continue
		}
		subtotal += Money(it.Qty) * it.UnitPriceCents
	}

	// discount = round(subtotal * discountPct / 100)
	discount := RoundPercent(subtotal, discountPct)

	// tax is applied on (subtotal - discount)
	taxBase := subtotal - discount
	if taxBase < 0 {
		taxBase = 0
	}
	tax := RoundPercent(taxBase, taxPct)

	total := subtotal - discount + tax
	return Totals{
		Subtotal: subtotal,
		Discount: discount,
		Tax:      tax,
		Total:    total,
	}
}

// roundPercent does (amount * pct / 100) with halves-up rounding.
func RoundPercent(amount Money, pct int) Money {
	if pct <= 0 || amount == 0 {
		return 0
	}
	// amount * pct may overflow if absurdly large; for typical invoices it's fine.
	// Add 50 for halves-up (since we're dividing by 100).
	raw := int64(amount) * int64(pct)
	return Money((raw + 50) / 100)
}
