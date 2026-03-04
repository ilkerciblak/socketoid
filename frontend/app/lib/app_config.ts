type AppConfig = {
  BASE_URL: string;
};

export const app_config: AppConfig = {
  BASE_URL: process.env.SERVER_URL!
} as AppConfig;
