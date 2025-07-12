# Skill Swap Backend - Heroku Deployment

This guide will help you deploy the Skill Swap backend to Heroku using Docker containers.

## Prerequisites

1. **Heroku CLI installed** - Download from [heroku.com](https://devcenter.heroku.com/articles/heroku-cli)
2. **Git repository initialized** - Your code should be in a Git repository
3. **Heroku account** - Sign up at [heroku.com](https://heroku.com)

## Quick Deployment

### Option 1: Automated Script
```bash
./deploy-heroku.sh
```

The script will:
- Create a new Heroku app (or use existing)
- Set up PostgreSQL database
- Configure environment variables
- Deploy using Docker container
- Set up SSL and production settings

### Option 2: Manual Steps

1. **Login to Heroku**
```bash
heroku login
```

2. **Create a new app**
```bash
heroku create your-app-name --region us
```

3. **Set stack to container**
```bash
heroku stack:set container -a your-app-name
```

4. **Add PostgreSQL addon**
```bash
heroku addons:create heroku-postgresql:essential-0 -a your-app-name
```

5. **Set environment variables**
```bash
# Generate secure JWT secret
JWT_SECRET=$(openssl rand -base64 32)
heroku config:set JWT_SECRET="$JWT_SECRET" -a your-app-name

# Set base URL
heroku config:set BASE_URL="https://your-app-name.herokuapp.com" -a your-app-name

# Set production mode
heroku config:set GIN_MODE=release -a your-app-name
```

6. **Deploy**
```bash
git add .
git commit -m "Deploy to Heroku"
heroku git:remote -a your-app-name
git push heroku main
```

## Environment Variables

The following environment variables are automatically configured:

- `DATABASE_URL` - PostgreSQL connection (auto-configured by Heroku)
- `PORT` - Application port (auto-configured by Heroku)
- `JWT_SECRET` - Secure JWT signing key
- `BASE_URL` - Your app's public URL
- `GIN_MODE` - Set to "release" for production

## API Endpoints

Once deployed, your API will be available at:

- **Base URL**: `https://your-app-name.herokuapp.com/api/v1`
- **Health Check**: `https://your-app-name.herokuapp.com/health`
- **API Documentation**: Available in the `backend_api_endpoints.md` file

## Post-Deployment

1. **Test the deployment**
```bash
curl https://your-app-name.herokuapp.com/health
```

2. **View logs**
```bash
heroku logs --tail -a your-app-name
```

3. **Open the app**
```bash
heroku open -a your-app-name
```

## Database Management

Your PostgreSQL database is automatically provisioned. To access it:

```bash
# Connect to database
heroku pg:psql -a your-app-name

# View database info
heroku pg:info -a your-app-name

# Reset database (⚠️ DESTRUCTIVE)
heroku pg:reset DATABASE_URL -a your-app-name --confirm your-app-name
```

## Monitoring and Scaling

```bash
# View app metrics
heroku ps -a your-app-name

# Scale dynos
heroku ps:scale web=2 -a your-app-name

# View configuration
heroku config -a your-app-name
```

## Troubleshooting

1. **Build fails**: Check Dockerfile and ensure all dependencies are correct
2. **Database connection issues**: Verify `DATABASE_URL` is set correctly
3. **Application errors**: Check logs with `heroku logs --tail`

## Security Notes

- JWT secrets are automatically generated and secured
- Database credentials are managed by Heroku
- HTTPS is enforced automatically
- CORS headers are configured for web requests

## Cost Information

- **Essential Postgres**: $5/month (recommended for production)
- **Eco Dyno**: $5/month (basic compute)
- **Total minimum**: ~$10/month for production deployment

For development/testing, you can use free tiers where available.

## Next Steps

After deployment:
1. Set up your frontend to use the new API URL
2. Configure domain name (optional)
3. Set up monitoring and alerting
4. Review and optimize performance
