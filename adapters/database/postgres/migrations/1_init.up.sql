-- +goose Up
CREATE TABLE devices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    jid VARCHAR(255) NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    is_active BOOLEAN NOT NULL DEFAULT FALSE
);

ALTER TABLE public.devices ADD CONSTRAINT devices_whatsmeow_device_fk FOREIGN KEY (jid) REFERENCES public.whatsmeow_device(jid) ON DELETE SET NULL ON UPDATE CASCADE;


-- +goose Down
DROP TABLE devices;