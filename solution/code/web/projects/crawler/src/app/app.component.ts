import { Component, OnInit } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { timeout } from 'rxjs/operators';

@Component({
  selector: 'crl-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnInit {
  title = 'crawler';
  sites: Observable<any>;
  url: string;
  tryUrls = ['/api', 'http://localhost:8080', 'http://localhost/api'];
  error: string;

  constructor(private http: HttpClient) {
  }

  ngOnInit() {
    this.tryUrl(0);
  }

  tryUrl(index) {
    const tryUrl = this.tryUrls[index];
    if (tryUrl !== undefined) {
      this.http.get(tryUrl, { responseType: 'text' })
        .pipe(timeout(100))
        .subscribe(
          response => {
            this.url = tryUrl;
            this.getSites();
          },
          error => {
            this.tryUrl(index + 1);
          }
        )
    } else {
      this.error = 'no running backend found'
    }
  }

  getSites() {
    this.sites = this.http.get(this.url + '/sites');
  }

  addSite(site) {
    this.http.post(this.url + '/sites', { url: site, interval: 100 }).subscribe(
      response => this.getSites(),
      error => this.error = error.message
    );
  }
}
