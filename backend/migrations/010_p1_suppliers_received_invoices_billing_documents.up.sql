-- P1: suppliers, received_invoices, billing_documents (idempotent for legacy + fresh installs).
BEGIN;

CREATE TABLE IF NOT EXISTS suppliers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(200) NOT NULL,
    contact_email VARCHAR(200),
    contact_phone VARCHAR(40),
    tax_id VARCHAR(64),
    notes TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_suppliers_deleted_at ON suppliers (deleted_at);
CREATE INDEX IF NOT EXISTS idx_suppliers_name ON suppliers (name);

CREATE TABLE IF NOT EXISTS received_invoices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    supplier_id UUID REFERENCES suppliers (id) ON DELETE SET NULL,
    vendor_name VARCHAR(200),
    category VARCHAR(80) NOT NULL,
    amount DOUBLE PRECISION NOT NULL,
    invoice_date TIMESTAMPTZ NOT NULL,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_received_invoices_supplier_id ON received_invoices (supplier_id);
CREATE INDEX IF NOT EXISTS idx_received_invoices_deleted_at ON received_invoices (deleted_at);
CREATE INDEX IF NOT EXISTS idx_received_invoices_invoice_date ON received_invoices (invoice_date);

CREATE TABLE IF NOT EXISTS billing_documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    kind VARCHAR(32) NOT NULL,
    title VARCHAR(200) NOT NULL,
    amount DOUBLE PRECISION NOT NULL,
    customer_id UUID,
    reference VARCHAR(120),
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_billing_documents_kind ON billing_documents (kind);
CREATE INDEX IF NOT EXISTS idx_billing_documents_customer_id ON billing_documents (customer_id);
CREATE INDEX IF NOT EXISTS idx_billing_documents_deleted_at ON billing_documents (deleted_at);

COMMIT;
