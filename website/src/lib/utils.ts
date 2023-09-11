import mailtoLink from "mailto-link";

export const mailFromOperator = mailtoLink({
  to: "contact@wheretopark.app",
  subject: "Contact from a Parking Operator",
  body: "Hello,\nI am a parking operator and I would like to add my parking lots to your app."
});

export const getRandomBetween = (min: number, max: number) => {
  return Math.floor(Math.random() * (max - min) + min);
};
