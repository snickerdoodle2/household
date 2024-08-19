export type Result<T, E> =
    | {
          isError: false;
          data: T;
      }
    | {
          isError: true;
          error: E;
      };
