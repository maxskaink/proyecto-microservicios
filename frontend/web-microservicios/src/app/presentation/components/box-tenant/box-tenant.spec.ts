import { ComponentFixture, TestBed } from '@angular/core/testing';

import { BoxTenant } from './box-tenant';

describe('BoxTenant', () => {
  let component: BoxTenant;
  let fixture: ComponentFixture<BoxTenant>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [BoxTenant]
    })
    .compileComponents();

    fixture = TestBed.createComponent(BoxTenant);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
