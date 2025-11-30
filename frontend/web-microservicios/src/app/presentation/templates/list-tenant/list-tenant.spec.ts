import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ListTenant } from './list-tenant';

describe('ListTenant', () => {
  let component: ListTenant;
  let fixture: ComponentFixture<ListTenant>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ListTenant]
    })
    .compileComponents();

    fixture = TestBed.createComponent(ListTenant);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
