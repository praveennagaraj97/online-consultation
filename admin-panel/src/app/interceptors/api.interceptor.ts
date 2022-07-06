import {
  HttpEvent,
  HttpHandler,
  HttpHeaders,
  HttpInterceptor,
  HttpRequest,
} from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { environment } from 'src/environments/environment';

@Injectable()
export class APiInterceptor implements HttpInterceptor {
  intercept(
    req: HttpRequest<any>,
    next: HttpHandler
  ): Observable<HttpEvent<any>> {
    const cloneReq = req.clone({
      url: `${environment.baseURL}${req.url}`,
      headers: new HttpHeaders({
        'Time-Zone': Intl.DateTimeFormat().resolvedOptions().timeZone,
      }),
      withCredentials: true,
    });

    return next.handle(cloneReq);
  }
}
