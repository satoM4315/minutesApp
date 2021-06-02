export type User = {
  id: number;
  name: string;
};

export type Message = {
  addedBy: User;
  id: number;
  message: string;
};

export type Meeting = {
  name: string;
  id: number;
};
