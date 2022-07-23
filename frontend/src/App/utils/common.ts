let lastCall = Date.now()
let lastCallTimer: ReturnType<typeof setTimeout>

export const debounce = (callee: Function, timeoutMs: number) => {

  return function perform() {

    let previousCall = lastCall;

    lastCall = Date.now();

    if (previousCall && lastCall - previousCall <= timeoutMs) {
      clearTimeout(lastCallTimer);
    }

    lastCallTimer = setTimeout(() => callee(), timeoutMs);

  };
}

export const isSameDate = function(date1:Date,date2:Date) {
  return (
    date1.getFullYear() === date2.getFullYear() &&
    date1.getMonth() === date2.getMonth() &&
    date1.getDate() === date2.getDate()
  );
}
