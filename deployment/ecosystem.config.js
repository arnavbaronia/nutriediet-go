// PM2 Ecosystem Configuration for Nutriediet New App
// Run with: pm2 start ecosystem.config.js

module.exports = {
  apps: [
    {
      name: 'nutriediet-go-api',
      script: './nutriediet-go',
      cwd: '/home/sk/mys/nutriediet-new/backend',
      instances: 1,
      exec_mode: 'fork',
      autorestart: true,
      watch: false,
      max_memory_restart: '500M',
      env: {
        PORT: '8080',
        GIN_MODE: 'release',
        NODE_ENV: 'production'
      },
      error_file: '/home/sk/mys/nutriediet-new/logs/go-api-error.log',
      out_file: '/home/sk/mys/nutriediet-new/logs/go-api-out.log',
      log_date_format: 'YYYY-MM-DD HH:mm:ss Z',
      merge_logs: true,
      min_uptime: '10s',
      max_restarts: 10,
      restart_delay: 4000
    }
  ]
};

// Note: React app is served as static files by Nginx, no PM2 process needed

