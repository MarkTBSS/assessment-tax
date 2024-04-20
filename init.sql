-- Create admin_settings table
CREATE TABLE IF NOT EXISTS public.admin_settings (
    id SERIAL PRIMARY KEY,
    personal_deduction DECIMAL(15, 2) NOT NULL,
    k_receipt DECIMAL(15, 2) NOT NULL
);

-- Populate admin_settings table with default values
INSERT INTO public.admin_settings (personal_deduction, k_receipt) VALUES (60000.00, 50000.00);