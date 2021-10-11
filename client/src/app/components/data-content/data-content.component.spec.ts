import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { DataContentComponent } from './data-content.component';

describe('DataContentComponent', () => {
  let component: DataContentComponent;
  let fixture: ComponentFixture<DataContentComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [ DataContentComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(DataContentComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
