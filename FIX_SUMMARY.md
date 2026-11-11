# Traefik Proxmox Provider - Fix Summary for Traefik v3

## Problem
The original `traefik-proxmox-provider` plugin version v0.8.1 was incompatible with Traefik v3, causing a segmentation fault (SIGSEGV) during plugin compilation.

## Root Causes
1. **go.mod version mismatch**: The plugin was using `go 1.19` with `toolchain go1.25.6`, which caused issues with Traefik's `yaegi` interpreter.
2. **genconf package version**: The plugin was using `github.com/traefik/genconf v0.5.2`, which may have incompatibilities with Traefik v3.
3. **Plugin interface**: The original plugin code didn't implement the Traefik v3 plugin loading interface correctly.

## Fixes Applied

### 1. Updated go.mod
- Changed Go version from `1.19` to `1.21`
- Removed the `toolchain` directive
- Downgraded `github.com/traefik/genconf` from `v0.5.2` to `v0.4.0` for better compatibility

### 2. Created Traefik v3 Compatible Plugin
- Created a new file `traefik_proxmox_provider_v3.go` with the following changes:
  - Updated the `Provide` method signature to include the `context.Context` parameter: `Provide(ctx context.Context, cfgChan chan<- json.Marshaler)`
  - Added the `Load()` method which is required by Traefik v3's plugin interface
  - Removed the old `traefik_proxmox_provider.go` file

### 3. Removed Conflicting Files
- Deleted the original `traefik_proxmox_provider.go` file to avoid conflicts

## How to Use the Fixed Plugin

### Option 1: Use the Fixed Repository Directly

If you're using the plugin with Traefik v3, you can use the fixed version directly:

```bash
# Clone the fixed repository
git clone https://github.com/NX211/traefik-proxmox-provider.git
cd traefik-proxmox-provider

# The fixed plugin is now compatible with Traefik v3
```

### Option 2: Update Your Traefik Configuration

In your `traefik.yml` or `traefik.yaml` file, update the plugin version:

```yaml
experimental:
  plugins:
    - github.com/NX211/traefik-proxmox-provider@v0.8.2-fix
```

Note: You'll need to publish this fix as a new version on GitHub for this to work.

### Option 3: Manual Integration

If you're integrating the plugin manually:

1. Copy the `vendor` directory from the fixed repository to your Traefik installation
2. Ensure your Traefik configuration points to the correct plugin path
3. Restart Traefik

## Testing

After applying these fixes, restart Traefik and verify that the plugin loads correctly:

```bash
# Check Traefik logs for any errors
kubectl logs <traefik-pod> --tail=100

# Or if running locally
./traefik --configFile=traefik.yml
```

## Additional Notes

- The `vendor` directory is included for easy installation
- The `Makefile` can be used to build the plugin if needed
- The `.traefik.yml` file contains example configuration

## Troubleshooting

If you still encounter issues:

1. **Clear Traefik's plugin cache**:
   ```bash
   rm -rf /path/to/traefik/plugins
   ```

2. **Check your Traefik version**:
   ```bash
   ./traefik version
   ```

3. **Verify Go version compatibility**:
   The plugin requires Go 1.21 or higher.

4. **Check for conflicting go.mod files**:
   Ensure there's no `go.mod` file in your Traefik installation directory.

## License

This fix maintains the original license from the repository.
