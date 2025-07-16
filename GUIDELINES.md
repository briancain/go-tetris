# Tetris Guidelines Implementation

This document outlines how this project implements the official Tetris Guidelines.

## Indispensable Rules

### Playfield Dimensions
- ✅ 10 cells wide and 22 cells tall (with top 2 rows hidden)
- ✅ Standard playfield dimensions with proper game over detection

### Tetrimino Colors
- ✅ Cyan I
- ✅ Yellow O
- ✅ Purple T
- ✅ Green S
- ✅ Red Z
- ✅ Blue J
- ✅ Orange L

### Tetrimino Start Locations
- ✅ I and O spawn in the middle columns
- ✅ J, L, S, T, Z spawn in the left-middle columns
- ✅ Tetrominoes spawn horizontally with J, L and T spawning flat-side first

### Rotation System
- ✅ Super Rotation System (SRS) with proper wall kicks
- ✅ Different wall kick data for I piece vs. other pieces

### Random Generator
- ✅ 7-bag randomizer ensuring all 7 pieces appear exactly once before any repeats

### Hold Piece
- ✅ Player can hold a piece for later use
- ✅ Hold cannot be used again until after the piece locks down

### Ghost Piece
- ✅ Shows where the current piece will land

### Terminology
- ✅ Uses "Tetriminos" as per guidelines
- ✅ Uses proper piece names (I, J, L, O, S, T, Z)

### Branding
- ✅ Support for Roger Dean's Tetris logo (requires user to provide the actual logo file)

## Scoring and Mechanics

### T-Spin Detection
- ✅ Implements 3-corner T rule for T-spin detection
- ✅ Bonus scoring for T-spins

### Back-to-Back Bonus
- ✅ 50% bonus for consecutive special clears (Tetris or T-spin)

### Leveling System
- ✅ Player levels up by clearing lines
- ✅ Speed increases with level

## Future Enhancements

### Audio
- [ ] Add Korobeiniki (Tetris theme song)
- [ ] Consider adding Katjusha or Kalinka songs

## References

- [Tetris Guidelines](https://tetris.fandom.com/wiki/Tetris_Guideline)
- [Super Rotation System](https://tetris.fandom.com/wiki/SRS)
- [7-Bag Random Generator](https://tetris.fandom.com/wiki/Random_Generator)
