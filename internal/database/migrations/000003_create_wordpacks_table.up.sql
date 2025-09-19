-- Create wordpacks table (must come before games due to foreign key)
CREATE TABLE wordpacks (
    id SERIAL PRIMARY KEY,
    name VARCHAR(256) UNIQUE NOT NULL,
    description VARCHAR(300),
    created_by UUID,
    is_default BOOLEAN DEFAULT FALSE,
    words TEXT[] NOT NULL,

    -- Foreign key constraint
    CONSTRAINT fk_wordpacks_created_by FOREIGN KEY (created_by) REFERENCES users(id)
);

CREATE INDEX idx_wordpacks_created_by ON wordpacks(created_by);
CREATE INDEX idx_wordpacks_is_default ON wordpacks(is_default);
CREATE INDEX idx_wordpacks_name ON wordpacks(name);


-- Insert some default word packs (optional) FIXME: remove this and seed properly
INSERT INTO wordpacks (name, description, is_default, words) VALUES
('Basic Words', 'Common everyday words', true, ARRAY['blue', 'bright', 'calm', 'cold', 'fast', 'funny', 'green', 'happy', 'loud', 'quick', 'sharp', 'small', 'soft', 'strong', 'tall', 'warm', 'young', 'build', 'chase', 'climb', 'dance', 'draw', 'eat', 'fly', 'grow', 'jump', 'kick', 'paint', 'run', 'sing', 'swim', 'talk', 'walk', 'write', 'yell', 'apple', 'bridge', 'cloud', 'desk', 'eagle', 'forest', 'glove', 'hammer', 'island', 'jacket', 'kitten', 'mirror', 'ocean', 'pencil', 'river', 'stone', 'table', 'window', 'badly', 'boldly', 'calmly', 'clearly', 'fast', 'gently', 'happily', 'loudly', 'neatly', 'often', 'quickly', 'rarely', 'silently', 'slowly', 'softly', 'tightly', 'well'])
ON CONFLICT (name) DO NOTHING;
