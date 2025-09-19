-- Create custom enum type
CREATE TYPE game_status AS ENUM ('Initial', 'Started', 'Finished');

-- Create games table
CREATE TABLE games (
    id UUID PRIMARY KEY NOT NULL,
    host_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    started_at TIMESTAMP,
    status game_status NOT NULL DEFAULT 'Initial',
    word_pack_id INTEGER NOT NULL,
    game_state JSONB NOT NULL,

    -- Foreign key constraints
    CONSTRAINT fk_games_host FOREIGN KEY (host_id) REFERENCES users(id),
    CONSTRAINT fk_games_word_pack FOREIGN KEY (word_pack_id) REFERENCES wordpacks(id)
);

-- Create indexes for better performance
CREATE INDEX idx_games_host_id ON games(host_id);
CREATE INDEX idx_games_status ON games(status);
CREATE INDEX idx_games_created_at ON games(created_at);
CREATE INDEX idx_games_word_pack_id ON games(word_pack_id);