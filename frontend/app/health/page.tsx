import { server } from "../lib/fetcher";

export default function Health() {
  type Health = {
    status: number;
    timee: string;
  };
  const HealthStatus = async () => {
    const healthData = (await server.GET<Health>("health")) as Health;

    return <h1>{healthData.status}</h1>;
  };

  return (
    <main>
      <HealthStatus />
    </main>
  );
}
