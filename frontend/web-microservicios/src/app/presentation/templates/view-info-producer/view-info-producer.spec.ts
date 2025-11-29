import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ViewInfoProducer } from './view-info-producer';

describe('ViewInfoProducer', () => {
  let component: ViewInfoProducer;
  let fixture: ComponentFixture<ViewInfoProducer>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ViewInfoProducer]
    })
    .compileComponents();

    fixture = TestBed.createComponent(ViewInfoProducer);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
