type AppConfig = {
  BASE_URL: string;
};

export const app_config: AppConfig = {
  BASE_URL: process.env.NEXT_PUBLIC_SERVER_URL!
} as AppConfig;
