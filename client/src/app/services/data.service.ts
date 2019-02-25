import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { APIResponse } from '../models/APIResponse'

const httpOptions = {
  headers: new HttpHeaders({
    'Content-Type':  'application/json'
  })
};

@Injectable({
  providedIn: 'root'
})
export class DataService {

  // Go server
  serverURL = "http://localhost:8080"

  constructor(
    private http : HttpClient
  ) { }

  // Load all data
  getAllData():Observable<APIResponse> {
    return this.http.get<APIResponse>(this.serverURL+'/getAllData')
  }

  // Write data
  writeData(data:any):Observable<APIResponse> {
    return this.http.post<APIResponse>(this.serverURL+'/writeData',data,httpOptions)
  }
}
