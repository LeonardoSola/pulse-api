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

ALTER TABLE devices ADD CONSTRAINT device_whatsmeow_device_fk FOREIGN KEY (jid) REFERENCES whatsmeow_device(jid) ON DELETE SET NULL ON UPDATE CASCADE;
CREATE INDEX idx_devices_jid ON devices (jid) WHERE deleted_at IS NULL;
CREATE INDEX idx_devices_is_active ON devices (is_active) WHERE deleted_at IS NULL;
CREATE INDEX idx_devices_created_at ON devices (created_at) WHERE deleted_at IS NULL;

CREATE TABLE messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	external_id VARCHAR(255) NOT NULL,
    device_id UUID NOT NULL,
    text TEXT NOT NULL,
    recipient VARCHAR(255) NOT NULL,
    recipient_type VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    FOREIGN KEY (device_id) REFERENCES devices(id) ON DELETE CASCADE ON UPDATE CASCADE
);
CREATE INDEX idx_messages_device_id ON messages (device_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_messages_recipient ON messages (recipient, recipient_type) WHERE deleted_at IS NULL;
CREATE INDEX idx_messages_device_created_at ON messages (device_id, created_at DESC) WHERE deleted_at IS NULL;

-- +goose Down
DROP TABLE devices;
DROP TABLE messages;