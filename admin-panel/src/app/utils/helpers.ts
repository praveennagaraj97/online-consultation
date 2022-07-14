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

export function convertFileUploadToBlobUrl(file: File) {
  if (file) {
    return URL.createObjectURL(file);
  }

  return '/assets/img-placeholder.png';
}
