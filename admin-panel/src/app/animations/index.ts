import { animate, animation, style } from '@angular/animations';

export const fadeInTransformAnimation = (time?: number) =>
  animation([
    // Initaial
    style({ opacity: 0 }),
    // Animate
    animate(time || 600, style({ opacity: 1 })),
  ]);
