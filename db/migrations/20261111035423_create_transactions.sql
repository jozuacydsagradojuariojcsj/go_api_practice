-- +goose Up
CREATE TABLE transactions (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    from_user_id UUID REFERENCES users(id),  -- NULL = deposit
    to_user_id   UUID REFERENCES users(id),  -- NULL = withdrawal
    amount       BIGINT NOT NULL CHECK (amount > 0),
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    CHECK (from_user_id IS NOT NULL OR to_user_id IS NOT NULL)
);

-- +goose Down
DROP TABLE transactions;