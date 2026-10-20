import { NotFound } from "~/components/not-found"

export function meta() {
  return [{ title: "Not found | Onevra" }]
}

export default function NotFoundRoute() {
  return <NotFound />
}
