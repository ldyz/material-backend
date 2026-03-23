-- Migration: Fix inbound_items table id column
-- Description: Add auto-increment sequence to existing inbound_items table

-- Check if table exists and fix the id column
DO $$
BEGIN
    -- Check if the table has a sequence, if not create one and set it as default
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'inbound_items') THEN
        -- Create sequence if it doesn't exist
        IF NOT EXISTS (SELECT 1 FROM pg_sequences WHERE schemaname = 'public' AND sequencename = 'inbound_items_id_seq') THEN
            CREATE SEQUENCE inbound_items_id_seq START 1;
        END IF;

        -- Set the id column to use the sequence
        EXECUTE format('ALTER TABLE inbound_items ALTER COLUMN id SET DEFAULT nextval(%L)', 'inbound_items_id_seq');

        -- Set existing id values to use the sequence if they are 0 or NULL
        EXECUTE 'ALTER TABLE inbound_items ALTER COLUMN id DROP DEFAULT';
        EXECUTE 'ALTER TABLE inbound_items ALTER COLUMN id DROP NOT NULL';
        EXECUTE 'ALTER TABLE inbound_items ALTER COLUMN id SET NOT NULL';
        EXECUTE format('ALTER TABLE inbound_items ALTER COLUMN id SET DEFAULT nextval(%L)', 'inbound_items_id_seq');
    END IF;
END $$;

-- Make stock_id nullable
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'inbound_items' AND column_name = 'stock_id' AND is_nullable = 'NO') THEN
        ALTER TABLE inbound_items ALTER COLUMN stock_id DROP NOT NULL;
    END IF;
END $$;
