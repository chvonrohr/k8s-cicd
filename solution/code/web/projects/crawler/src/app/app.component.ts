import { Component, OnInit } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

@Component({
  selector: 'crl-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnInit {
  title = 'crawler';
  sites: Observable<any>;

  constructor(private http: HttpClient) {
  }

  ngOnInit() {
    this.sites = this.http.get('http://localhost:8080/sites');
  }
}
