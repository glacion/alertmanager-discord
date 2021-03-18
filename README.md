# alertmanager-discord

Give this a webhook (with the DISCORD_WEBHOOK environment variable) and point it as a webhook on alertmanager, and it will post your alerts into a discord channel for you as they trigger:

Example alert manager config:

```yaml
receivers:
- name: 'discord_webhook'
  webhook_configs:
  - url: 'http://localhost:9094'
```

