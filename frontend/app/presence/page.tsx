"use client";
import { app_config } from "../lib/app_config";
import useSSE from "../lib/hooks/use_sse";
import { useActionState, useState } from "react";
import OnlineUser from "./_components/online_user";

export default function Presence() {
  const [userList, setUserList] = useState<{ name: string; user_id: string }[]>(
    [],
  );

  const { connect } = useSSE(app_config.PUBLIC_API + "/events", {
    auto: false,
    handlers: [
      {
        type: "presence.init",
        handler(event) {
          try {
            const users = JSON.parse(event.data);

            setUserList(users || []);
          } catch (error) {
            console.log(error);
          }
        },
      },
      {
        type: "user.joined",
        handler(event) {
          const user = JSON.parse(event.data);
          setUserList((prev) => [...prev, user]);
        },
      },
      {
        type: "user.left",
        handler(event) {
          const user = JSON.parse(event.data) as {
            name: string;
            user_id: string;
          };
          setUserList((prev) => prev.filter((u) => u.user_id != user.user_id));
        },
      },
    ],
  });
  const submitName = (prevState: any, formData: FormData) => {
    const firstName: string | null = formData.get("firstname") as string | null;
    connect({ key: "name", val: firstName ?? "" });
    return {
      enteredValues: { ...prevState, firstName },
    };
  };

  //const [state, formAction] = useActionState(submitName, {
  //    enteredValues: {},
  // });
  //

  const [state, formAction] = useActionState(submitName, { enteredValues: {} });
  return (
    <>
      <div className="flex flex-col justify-center items-center h-dvh">
        <form
          action={formAction}
          className="flex flex-col bg-white rounded-md p-5 gap-y-3"
        >
          <label
            htmlFor="firstname"
            className="block w-full font-mono text-gray-500"
          >
            Enter your firstname to connect
          </label>
          <div className="w-full flex flex-row justify-between">
            <input
              type="text"
              name="firstname"
              id="firstname"
              placeholder="Uvuvwevwevwe Onyetenyevwe Ugwemuhwem Osas"
              defaultValue={state.enteredValues?.firstName}
              className="border border-gray-200 rounded-sm px-1 text-gray-900 placeholder-gray-500"
            />
            <button
              className="rounded-xl bg-blue-900 hover:bg-blue-500 cursor-pointer  px-1.5 py-1.5 font-bold"
              type="submit"
            >
              Connect
            </button>
          </div>
        </form>
        <ul className="flex flex-col gap-y-2 items-start my-3">
          <h3 className="font-bold text-gray-500">Online Users</h3>
          {userList?.length > 0 &&
            userList.map((u) => (
              <OnlineUser username={u.name} id={u.user_id} key={u.user_id} />
            ))}
        </ul>
      </div>
    </>
  );
}
