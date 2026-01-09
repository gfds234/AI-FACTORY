# Prompt Improvements Test

## Test Case: Verify Template Detection Works

This file tests that our template detection catches README placeholders.

### Test Input (Simulated Bad LLM Output):

```
## Project Title
Project description.

## Technologies Used
- Technology 1
- Technology 2

### Prerequisites
What things you need to install the software and how to install them:
```
Give examples
```

### Installing
A step by step series of examples that tell you how to get a development env running:
```
Give the example command
```
```

### Expected Behavior:
- Template detection should catch "Give examples"
- Files should NOT be saved to projects/
- Artifact path should indicate "(no files extracted - check artifact for template errors)"

### Test Verification:
Run the code generation and check logs for:
- `[ERROR] Detected README template in file`
- `[ERROR] Code generation produced README template instead of actual code`
