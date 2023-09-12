import mailtoLink from "mailto-link";
import type { Thing, WithContext } from "schema-dts";

export const mailFromOperator = mailtoLink({
  to: "contact@wheretopark.app",
  subject: "Contact from a Parking Operator",
  body: "Hello,\nI am a parking operator and I would like to add my parking lots to your app."
});

export const getRandomBetween = (min: number, max: number) => {
  return Math.floor(Math.random() * (max - min) + min);
};

export const serializeSchema = <T extends Thing>(thing: WithContext<T>): string => {
  return `<script type="application/ld+json">${JSON.stringify(
    thing,
    null,
    2
  )}</script>`
}