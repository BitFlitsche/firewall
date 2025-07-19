-- Seed data for charset_rules table
-- This includes major Unicode scripts and their associated languages
-- Status is set to 'allowed' by default, can be changed as needed

INSERT INTO charset_rules (charset, status, created_at, updated_at) VALUES
-- Basic ASCII and Latin scripts
('ASCII', 'allowed', NOW(), NOW()),
('Latin', 'allowed', NOW(), NOW()),
('Vietnamese', 'allowed', NOW(), NOW()),

-- Cyrillic scripts (Russian, Ukrainian, Bulgarian, Serbian, etc.)
('Cyrillic', 'allowed', NOW(), NOW()),

-- Arabic scripts (Arabic, Persian, Urdu, etc.)
('Arabic', 'allowed', NOW(), NOW()),

-- Hebrew script
('Hebrew', 'allowed', NOW(), NOW()),

-- Greek script
('Greek', 'allowed', NOW(), NOW()),

-- South Asian scripts
('Devanagari', 'allowed', NOW(), NOW()),  -- Hindi, Sanskrit, Marathi, etc.
('Bengali', 'allowed', NOW(), NOW()),     -- Bengali, Assamese
('Tamil', 'allowed', NOW(), NOW()),       -- Tamil
('Telugu', 'allowed', NOW(), NOW()),      -- Telugu
('Kannada', 'allowed', NOW(), NOW()),     -- Kannada
('Malayalam', 'allowed', NOW(), NOW()),   -- Malayalam
('Gujarati', 'allowed', NOW(), NOW()),    -- Gujarati
('Gurmukhi', 'allowed', NOW(), NOW()),    -- Punjabi
('Oriya', 'allowed', NOW(), NOW()),       -- Odia
('Sinhala', 'allowed', NOW(), NOW()),     -- Sinhala

-- Southeast Asian scripts
('Thai', 'allowed', NOW(), NOW()),        -- Thai
('Lao', 'allowed', NOW(), NOW()),         -- Lao
('Khmer', 'allowed', NOW(), NOW()),       -- Khmer
('Myanmar', 'allowed', NOW(), NOW()),     -- Burmese

-- East Asian scripts
('Chinese', 'allowed', NOW(), NOW()),     -- Chinese (Simplified & Traditional)
('Japanese', 'allowed', NOW(), NOW()),    -- Japanese (Hiragana, Katakana, Kanji)
('Korean', 'allowed', NOW(), NOW()),      -- Korean (Hangul)

-- Other scripts
('Armenian', 'allowed', NOW(), NOW()),    -- Armenian
('Georgian', 'allowed', NOW(), NOW()),    -- Georgian
('Ethiopic', 'allowed', NOW(), NOW()),    -- Amharic, Tigrinya, etc.
('Mongolian', 'allowed', NOW(), NOW()),   -- Mongolian
('Tibetan', 'allowed', NOW(), NOW()),     -- Tibetan

-- Special categories
('Mixed', 'allowed', NOW(), NOW()),       -- Mixed scripts
('UTF-8', 'allowed', NOW(), NOW()),       -- UTF-8 encoded text
('Other', 'allowed', NOW(), NOW());       -- Other unrecognized scripts 