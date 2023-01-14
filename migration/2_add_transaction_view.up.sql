CREATE MATERIALIZED VIEW IF NOT EXISTS transaction_query AS
SELECT *
FROM (SELECT id,
             transaction_id,
             status,
             amount,
             product_code,
             user_id,
             created_at,
             ROW_NUMBER() OVER (
                 PARTITION BY transaction_id
                 ORDER BY created_at DESC
                 )
      FROM transaction) transaction_temp
WHERE row_number = 1;
