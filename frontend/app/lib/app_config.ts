type AppConfig = {
  BASE_URL: string;
  PUBLIC_API: string;
};

export const app_config: AppConfig = {
  BASE_URL: process.env.API_URL!,
  PUBLIC_API: process.env.NEXT_PUBLIC_SERVER_URL!,
} as AppConfig;
