import { NgModule } from '@angular/core';
import { ReplaceUnderScorePipe } from './replace-underscore.pipe';

@NgModule({
  declarations: [ReplaceUnderScorePipe],
  exports: [ReplaceUnderScorePipe],
})
export class PipesModule {}
