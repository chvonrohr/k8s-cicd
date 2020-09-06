import { Component, OnInit } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import {shareReplay, tap, timeout} from 'rxjs/operators';
import {Crawl, Page, Site} from './model';

@Component({
  selector: 'crl-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnInit {
  title = 'Crawler';
  sites: Observable<Site[]>;
  url: string;
  selectedSite: Site;
  crawls: Observable<Crawl[]>;
  pages: Observable<Page[]>;
  tryUrls = [
    'http://' + window.location.hostname,
    'http://' + window.location.hostname + ':8080',
    '/api',
    'http://localhost:8080',
    'http://localhost/api'];
  error: string;
  selectedCrawl: Crawl;

  constructor(private http: HttpClient) {
  }

  ngOnInit(): void {
    this.tryUrl(0);
  }

  tryUrl(index): void {
    const tryUrl = this.tryUrls[index];
    if (tryUrl !== undefined) {
      this.http.get(tryUrl, { responseType: 'text' })
        .pipe(timeout(100))
        .subscribe(
          response => {
            if (response === 'backend works') {
              this.url = tryUrl;
              this.getSites();
            } else {
              this.tryUrl(index + 1);
            }
          },
          error => {
            this.tryUrl(index + 1);
          }
        );
    } else {
      this.error = 'no running backend found';
    }
  }

  getSites(): void {
    this.sites = this.http.get<Site[]>(this.url + '/sites');
  }

  addSite(site): void {
    // todo: @wingsuitist - interval is 100 milliseconds
    this.http.post(this.url + '/sites', { url: site, interval: 100 }).subscribe(
      response => this.getSites(),
      error => this.error = error.message
    );
  }

  selectSite(site: Site): void {
    this.selectedSite = site;
    this.loadCrawlsForSite(site);
  }

  private loadCrawlsForSite(site: Site): void {
    this.crawls = this.http.get<Crawl[]>(`${this.url}/crawls?site=${site.ID}`)
      .pipe(
        shareReplay(1),
        tap(crawls => {
          if (crawls.length > 0 ) {
            this.selectCrawl(crawls[crawls.length - 1]);
          }
        })
      );
    this.crawls.subscribe();
  }

  selectCrawl(crawl: Crawl): void {
    this.selectedCrawl = crawl;
    this.loadPagesForCrawl(crawl);
  }

  private loadPagesForCrawl(crawl: Crawl): void {
    this.pages = this.http.get<Page[]>(`${this.url}/pages?site=${crawl.ID}`);
  }

}
