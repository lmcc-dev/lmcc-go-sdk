<!--
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This document was collaboratively developed by Martin and AI Assistant.
-->

## Migration from Standard Library `errors`

Migrating from Go's standard library `errors` package to `pkg/errors` can be done incrementally. The `pkg/errors` module is designed to be largely compatible, especially with error wrapping and checking via `standardErrors.Is` and `standardErrors.As`.

Key benefits of migrating include:
- **Stack Traces**: Automatic stack trace capture on error creation.
- **Error Codes (`Coder`)**: Standardized way to categorize errors and associate them with HTTP statuses or other metadata.
- **Structured Formatting**: Richer error messages with `%+v`.
- **Error Aggregation**: `ErrorGroup` for collecting multiple errors. 