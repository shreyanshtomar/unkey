import { auth } from "@clerk/nextjs";
import { redirect } from "next/navigation";

export default function AuthLayout(props: { children: React.ReactNode }) {
  const { userId } = auth();

  if (userId) {
    return redirect("/app/apis");
  }
  return (
    <>
      <div className="grid h-screen grid-cols-1 bg-white place-items-center">
        <div className="container">{props.children}</div>
      </div>
    </>
  );
}
