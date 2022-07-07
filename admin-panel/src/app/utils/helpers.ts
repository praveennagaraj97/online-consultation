import { Subscription } from 'rxjs';

export function clearSubscriptions(subs: Subscription[]) {
  if (!subs.length) {
    return;
  }

  subs.forEach((sub) => {
    if (sub instanceof Subscription) {
      sub.unsubscribe();
    }
  });
}
