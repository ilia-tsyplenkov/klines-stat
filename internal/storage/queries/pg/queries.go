package pg

const (
	InsertKLineQuery = `
    INSERT INTO trades.klines (pair, time_frame, open_price, high_price, low_price, close_price, utc_begin, utc_end)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `

	InsertRecentTradeQuery = `
    INSERT INTO trades.recent_trade (tid, pair, price, amount, side, ts)
    VALUES ($1, $2, $3, $4, $5, $6)
    `
)
